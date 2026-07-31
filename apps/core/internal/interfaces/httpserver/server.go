package httpserver

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

const shutdownTimeout = 10 * time.Second

type HealthChecker interface {
	Ping(context.Context) error
	SchemaVersion() int64
}

type DirectoryPicker interface {
	Select(context.Context) (path string, selected bool, err error)
}

type Server struct {
	echo            *echo.Echo
	health          HealthChecker
	directoryPicker DirectoryPicker
	projects        ProjectService
	logger          *zap.Logger
}

type Dependencies struct {
	Health          HealthChecker
	DirectoryPicker DirectoryPicker
	Projects        ProjectService
	Logger          *zap.Logger
	WebAssets       fs.FS
}

func New(dependencies Dependencies) (*Server, error) {
	if dependencies.Health == nil {
		return nil, errors.New("health checker is required")
	}
	if dependencies.Projects == nil {
		return nil, errors.New("project service is required")
	}
	if dependencies.DirectoryPicker == nil {
		return nil, errors.New("directory picker is required")
	}
	if dependencies.Logger == nil {
		dependencies.Logger = zap.NewNop()
	}

	e := echo.New()
	e.Use(middleware.Recover())
	server := &Server{
		echo:            e,
		health:          dependencies.Health,
		directoryPicker: dependencies.DirectoryPicker,
		projects:        dependencies.Projects,
		logger:          dependencies.Logger,
	}
	e.GET("/healthz", server.healthHandler)
	e.POST("/api/v1/system/select-directory", server.selectDirectoryHandler)
	e.GET("/api/v1/projects", server.listProjectsHandler)
	e.POST("/api/v1/projects", server.createProjectHandler)
	e.GET("/api/v1/projects/:project_id", server.getProjectHandler)
	if dependencies.WebAssets != nil {
		if _, err := fs.Stat(dependencies.WebAssets, "index.html"); err != nil {
			return nil, fmt.Errorf("embedded Web UI is missing index.html: %w", err)
		}
		e.GET("/*", echo.WrapHandler(spaHandler(dependencies.WebAssets)))
	}

	return server, nil
}

func (s *Server) Handler() http.Handler {
	return s.echo
}

func (s *Server) Run(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", address, err)
	}
	return s.Serve(ctx, listener)
}

func (s *Server) Serve(ctx context.Context, listener net.Listener) error {
	httpServer := &http.Server{
		Handler:           s.echo,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- httpServer.Serve(listener)
	}()

	s.logger.Info("Core HTTP server listening", zap.String("address", listener.Addr().String()))
	select {
	case err := <-serveErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve Core HTTP: %w", err)
	case <-ctx.Done():
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	shutdownErr := httpServer.Shutdown(shutdownContext)
	serveErr := <-serveErrors
	if errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}
	if shutdownErr != nil {
		shutdownErr = fmt.Errorf("gracefully stop Core HTTP: %w", shutdownErr)
	}
	if serveErr != nil {
		serveErr = fmt.Errorf("serve Core HTTP during shutdown: %w", serveErr)
	}
	return errors.Join(shutdownErr, serveErr)
}

func (s *Server) healthHandler(c *echo.Context) error {
	if err := s.health.Ping(c.Request().Context()); err != nil {
		s.logger.Error("Core health check failed", zap.Error(err))
		return c.JSON(http.StatusServiceUnavailable, map[string]any{
			"status":         "unavailable",
			"database":       "unavailable",
			"schema_version": s.health.SchemaVersion(),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"status":         "ok",
		"database":       "ok",
		"schema_version": s.health.SchemaVersion(),
	})
}

func spaHandler(webAssets fs.FS) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		assetPath := strings.TrimPrefix(path.Clean("/"+request.URL.Path), "/")
		if assetPath == "" || assetPath == "." {
			assetPath = "index.html"
		}
		info, err := fs.Stat(webAssets, assetPath)
		if err != nil || info.IsDir() {
			assetPath = "index.html"
		}

		content, err := fs.ReadFile(webAssets, assetPath)
		if err != nil {
			http.Error(response, "Web UI asset unavailable", http.StatusInternalServerError)
			return
		}
		if contentType := mime.TypeByExtension(path.Ext(assetPath)); contentType != "" {
			response.Header().Set("Content-Type", contentType)
		}
		if assetPath == "index.html" {
			response.Header().Set("Cache-Control", "no-cache")
		} else {
			response.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		response.WriteHeader(http.StatusOK)
		_, _ = response.Write(content)
	})
}
