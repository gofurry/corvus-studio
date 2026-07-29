package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	clearConfigEnvironment(t)

	root := t.TempDir()
	cfg, err := load(LoadOptions{}, root)
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}

	wantDataDir := filepath.Join(root, "corvus-studio")
	if cfg.Runtime.DataDir != wantDataDir {
		t.Fatalf("data dir = %q, want %q", cfg.Runtime.DataDir, wantDataDir)
	}
	if cfg.Server.Host != DefaultHost || cfg.Server.Port != DefaultPort {
		t.Fatalf("server = %s, want %s:%d", cfg.Server.Address(), DefaultHost, DefaultPort)
	}
	if cfg.Storage.Path != filepath.Join(wantDataDir, "corvus.db") {
		t.Fatalf("storage path = %q", cfg.Storage.Path)
	}
	if cfg.Logging.Path != filepath.Join(wantDataDir, "logs", "corvus.log") {
		t.Fatalf("log path = %q", cfg.Logging.Path)
	}
	if cfg.LoadedFile != "" {
		t.Fatalf("loaded file = %q, want none", cfg.LoadedFile)
	}
}

func TestLoadPrecedence(t *testing.T) {
	clearConfigEnvironment(t)

	root := t.TempDir()
	configPath := filepath.Join(root, "config.yaml")
	content := []byte("server:\n  port: 9100\nlogging:\n  level: debug\n")
	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("CORVUS_SERVER_PORT", "9200")
	t.Setenv("CORVUS_LOGGING_LEVEL", "warn")

	cliPort := 9300
	cliLevel := "error"
	cfg, err := load(LoadOptions{
		ConfigFile: configPath,
		Overrides: Overrides{
			Port:     &cliPort,
			LogLevel: &cliLevel,
		},
	}, root)
	if err != nil {
		t.Fatalf("load precedence config: %v", err)
	}

	if cfg.Server.Port != 9300 {
		t.Fatalf("port = %d, want CLI value 9300", cfg.Server.Port)
	}
	if cfg.Logging.Level != "error" {
		t.Fatalf("log level = %q, want CLI value error", cfg.Logging.Level)
	}
	if cfg.LoadedFile != configPath {
		t.Fatalf("loaded file = %q, want %q", cfg.LoadedFile, configPath)
	}
}

func TestLoadEnvironmentOverridesFile(t *testing.T) {
	clearConfigEnvironment(t)

	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	content := []byte(`{"server":{"port":9100},"logging":{"level":"debug"}}`)
	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("CORVUS_SERVER_PORT", "9200")
	t.Setenv("CORVUS_LOGGING_LEVEL", "warn")

	cfg, err := load(LoadOptions{ConfigFile: configPath}, root)
	if err != nil {
		t.Fatalf("load environment config: %v", err)
	}
	if cfg.Server.Port != 9200 || cfg.Logging.Level != "warn" {
		t.Fatalf("environment precedence not applied: port=%d level=%q", cfg.Server.Port, cfg.Logging.Level)
	}
}

func TestLoadRejectsUnsafeHost(t *testing.T) {
	clearConfigEnvironment(t)

	host := "0.0.0.0"
	_, err := load(LoadOptions{Overrides: Overrides{Host: &host}}, t.TempDir())
	if err == nil {
		t.Fatal("load unsafe host returned nil error")
	}
}

func TestLoadRejectsMissingExplicitConfig(t *testing.T) {
	clearConfigEnvironment(t)

	_, err := load(LoadOptions{ConfigFile: filepath.Join(t.TempDir(), "missing.yaml")}, t.TempDir())
	if err == nil {
		t.Fatal("load missing explicit config returned nil error")
	}
}

func clearConfigEnvironment(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"CORVUS_RUNTIME_DATA_DIR",
		"CORVUS_SERVER_HOST",
		"CORVUS_SERVER_PORT",
		"CORVUS_STORAGE_PATH",
		"CORVUS_LOGGING_PATH",
		"CORVUS_LOGGING_LEVEL",
		"CORVUS_LOGGING_MAX_SIZE_MB",
		"CORVUS_LOGGING_MAX_BACKUPS",
		"CORVUS_LOGGING_MAX_AGE_DAYS",
		"CORVUS_LOGGING_COMPRESS",
	} {
		t.Setenv(key, "")
	}
}
