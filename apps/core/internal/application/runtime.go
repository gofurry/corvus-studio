package application

import (
	"context"
	"errors"
	"fmt"

	projectapp "github.com/gofurry/corvus-studio/apps/core/internal/application/project"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/config"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/logging"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/projectpath"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/storage"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/httpserver"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/webui"
	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg config.Config) (runErr error) {
	logHandle, err := logging.New(cfg.Logging)
	if err != nil {
		return fmt.Errorf("initialize logging: %w", err)
	}
	defer func() {
		runErr = errors.Join(runErr, logHandle.Close())
	}()
	logger := logHandle.Logger

	logger.Info(
		"starting Corvus Core",
		zap.String("address", cfg.Server.Address()),
		zap.String("data_dir", cfg.Runtime.DataDir),
		zap.String("database", cfg.Storage.Path),
		zap.String("config_file", cfg.LoadedFile),
	)

	store, err := storage.Open(ctx, cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("initialize storage: %w", err)
	}
	defer func() {
		runErr = errors.Join(runErr, store.Close())
	}()
	if backupPath := store.LastBackupPath(); backupPath != "" {
		logger.Info("created pre-migration database backup", zap.String("path", backupPath))
	}
	logger.Info("database ready", zap.Int64("schema_version", store.SchemaVersion()))

	projectRepository, err := storage.NewProjectRepository(store)
	if err != nil {
		return fmt.Errorf("initialize Project repository: %w", err)
	}
	projectService, err := projectapp.NewService(
		projectRepository,
		projectpath.Resolver{},
		projectapp.UUIDv7Generator{},
		projectapp.SystemClock{},
	)
	if err != nil {
		return fmt.Errorf("initialize Project service: %w", err)
	}

	assets, err := webui.Files()
	if err != nil {
		return fmt.Errorf("load embedded Web UI: %w", err)
	}
	server, err := httpserver.New(httpserver.Dependencies{
		Health:    store,
		Projects:  projectService,
		Logger:    logger,
		WebAssets: assets,
	})
	if err != nil {
		return fmt.Errorf("initialize HTTP server: %w", err)
	}
	if err := server.Run(ctx, cfg.Server.Address()); err != nil {
		return err
	}
	logger.Info("Corvus Core stopped")
	return nil
}
