package cli

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/config"
)

func TestServePassesCLIOverridesToConfigLoader(t *testing.T) {
	dataDir := t.TempDir()
	var loadedOptions config.LoadOptions
	var runtimeConfig config.Config
	command := NewRootCommand(Dependencies{
		LoadConfig: func(options config.LoadOptions) (config.Config, error) {
			loadedOptions = options
			return config.Config{
				Runtime: config.RuntimeConfig{DataDir: dataDir},
				Server:  config.ServerConfig{Host: config.DefaultHost, Port: 9876},
			}, nil
		},
		Run: func(_ context.Context, cfg config.Config) error {
			runtimeConfig = cfg
			return nil
		},
	})
	configPath := filepath.Join(dataDir, "config.yaml")
	command.SetArgs([]string{
		"--config", configPath,
		"serve",
		"--data-dir", dataDir,
		"--host", "localhost",
		"--port", "9876",
		"--database", filepath.Join(dataDir, "custom.db"),
		"--log-file", filepath.Join(dataDir, "custom.log"),
		"--log-level", "debug",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("execute serve: %v", err)
	}
	if loadedOptions.ConfigFile != configPath {
		t.Fatalf("config file = %q, want %q", loadedOptions.ConfigFile, configPath)
	}
	if loadedOptions.Overrides.Port == nil || *loadedOptions.Overrides.Port != 9876 {
		t.Fatalf("port override = %#v", loadedOptions.Overrides.Port)
	}
	if loadedOptions.Overrides.Host == nil || *loadedOptions.Overrides.Host != "localhost" {
		t.Fatalf("host override = %#v", loadedOptions.Overrides.Host)
	}
	if loadedOptions.Overrides.DataDir == nil || *loadedOptions.Overrides.DataDir != dataDir {
		t.Fatalf("data-dir override = %#v", loadedOptions.Overrides.DataDir)
	}
	if runtimeConfig.Server.Port != 9876 {
		t.Fatalf("runtime config was not passed to runner: %#v", runtimeConfig)
	}
}

func TestServePropagatesCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var observed error
	command := NewRootCommand(Dependencies{
		LoadConfig: func(config.LoadOptions) (config.Config, error) {
			return config.Config{}, nil
		},
		Run: func(ctx context.Context, _ config.Config) error {
			observed = ctx.Err()
			return nil
		},
	})
	command.SetArgs([]string{"serve"})

	if err := command.ExecuteContext(ctx); err != nil {
		t.Fatalf("execute cancelled serve: %v", err)
	}
	if observed != context.Canceled {
		t.Fatalf("runner context error = %v, want context.Canceled", observed)
	}
}
