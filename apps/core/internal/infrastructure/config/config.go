package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultHost = "127.0.0.1"
	DefaultPort = 8765
)

var supportedConfigExtensions = []string{".json", ".yaml", ".yml"}

type Config struct {
	Runtime    RuntimeConfig `mapstructure:"runtime"`
	Server     ServerConfig  `mapstructure:"server"`
	Storage    StorageConfig `mapstructure:"storage"`
	Logging    LoggingConfig `mapstructure:"logging"`
	LoadedFile string        `mapstructure:"-"`
}

type RuntimeConfig struct {
	DataDir string `mapstructure:"data_dir"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (c ServerConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type StorageConfig struct {
	Path string `mapstructure:"path"`
}

type LoggingConfig struct {
	Path       string `mapstructure:"path"`
	Level      string `mapstructure:"level"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
}

type Overrides struct {
	DataDir    *string
	Host       *string
	Port       *int
	Storage    *string
	LogPath    *string
	LogLevel   *string
	MaxSizeMB  *int
	MaxBackups *int
	MaxAgeDays *int
	Compress   *bool
}

type LoadOptions struct {
	ConfigFile string
	Overrides  Overrides
}

func Load(options LoadOptions) (Config, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve user config directory: %w", err)
	}

	return load(options, userConfigDir)
}

func load(options LoadOptions, userConfigDir string) (Config, error) {
	v := viper.New()
	v.SetEnvPrefix("CORVUS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	defaultDataDir := filepath.Join(userConfigDir, "corvus-studio")
	v.SetDefault("runtime.data_dir", defaultDataDir)
	v.SetDefault("server.host", DefaultHost)
	v.SetDefault("server.port", DefaultPort)
	v.SetDefault("storage.path", "")
	v.SetDefault("logging.path", "")
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.max_size_mb", 10)
	v.SetDefault("logging.max_backups", 3)
	v.SetDefault("logging.max_age_days", 28)
	v.SetDefault("logging.compress", true)

	loadedFile, err := readConfig(v, options.ConfigFile, filepath.Join(userConfigDir, "corvus-studio"))
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode configuration: %w", err)
	}
	cfg.LoadedFile = loadedFile
	applyOverrides(&cfg, options.Overrides)

	if err := resolvePaths(&cfg); err != nil {
		return Config{}, err
	}
	if err := validate(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func readConfig(v *viper.Viper, explicitPath, searchDir string) (string, error) {
	if explicitPath != "" {
		extension := strings.ToLower(filepath.Ext(explicitPath))
		if !slices.Contains(supportedConfigExtensions, extension) {
			return "", fmt.Errorf("unsupported config extension %q: use YAML or JSON", extension)
		}

		absolutePath, err := filepath.Abs(explicitPath)
		if err != nil {
			return "", fmt.Errorf("resolve config path: %w", err)
		}
		v.SetConfigFile(absolutePath)
		if err := v.ReadInConfig(); err != nil {
			return "", fmt.Errorf("read config file %q: %w", absolutePath, err)
		}
		return v.ConfigFileUsed(), nil
	}

	for _, name := range []string{"config.yaml", "config.yml", "config.json"} {
		candidate := filepath.Join(searchDir, name)
		_, err := os.Stat(candidate)
		switch {
		case err == nil:
			v.SetConfigFile(candidate)
			if err := v.ReadInConfig(); err != nil {
				return "", fmt.Errorf("read config file %q: %w", candidate, err)
			}
			return v.ConfigFileUsed(), nil
		case errors.Is(err, os.ErrNotExist):
			continue
		default:
			return "", fmt.Errorf("inspect config file %q: %w", candidate, err)
		}
	}

	return "", nil
}

func applyOverrides(cfg *Config, overrides Overrides) {
	if overrides.DataDir != nil {
		cfg.Runtime.DataDir = *overrides.DataDir
	}
	if overrides.Host != nil {
		cfg.Server.Host = *overrides.Host
	}
	if overrides.Port != nil {
		cfg.Server.Port = *overrides.Port
	}
	if overrides.Storage != nil {
		cfg.Storage.Path = *overrides.Storage
	}
	if overrides.LogPath != nil {
		cfg.Logging.Path = *overrides.LogPath
	}
	if overrides.LogLevel != nil {
		cfg.Logging.Level = *overrides.LogLevel
	}
	if overrides.MaxSizeMB != nil {
		cfg.Logging.MaxSizeMB = *overrides.MaxSizeMB
	}
	if overrides.MaxBackups != nil {
		cfg.Logging.MaxBackups = *overrides.MaxBackups
	}
	if overrides.MaxAgeDays != nil {
		cfg.Logging.MaxAgeDays = *overrides.MaxAgeDays
	}
	if overrides.Compress != nil {
		cfg.Logging.Compress = *overrides.Compress
	}
}

func resolvePaths(cfg *Config) error {
	dataDir, err := absoluteCleanPath(cfg.Runtime.DataDir)
	if err != nil {
		return fmt.Errorf("resolve runtime data directory: %w", err)
	}
	cfg.Runtime.DataDir = dataDir

	if strings.TrimSpace(cfg.Storage.Path) == "" {
		cfg.Storage.Path = filepath.Join(dataDir, "corvus.db")
	} else {
		cfg.Storage.Path, err = absoluteCleanPath(cfg.Storage.Path)
		if err != nil {
			return fmt.Errorf("resolve database path: %w", err)
		}
	}

	if strings.TrimSpace(cfg.Logging.Path) == "" {
		cfg.Logging.Path = filepath.Join(dataDir, "logs", "corvus.log")
	} else {
		cfg.Logging.Path, err = absoluteCleanPath(cfg.Logging.Path)
		if err != nil {
			return fmt.Errorf("resolve log path: %w", err)
		}
	}

	return nil
}

func absoluteCleanPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("path is empty")
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(absolutePath), nil
}

func validate(cfg *Config) error {
	host := strings.TrimSpace(cfg.Server.Host)
	if host != "localhost" {
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return fmt.Errorf("server host %q is not loopback; Phase 1 does not expose an unauthenticated Core", cfg.Server.Host)
		}
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server port %d is outside 1..65535", cfg.Server.Port)
	}

	level := strings.ToLower(strings.TrimSpace(cfg.Logging.Level))
	if !slices.Contains([]string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}, level) {
		return fmt.Errorf("unsupported log level %q", cfg.Logging.Level)
	}
	cfg.Logging.Level = level
	if cfg.Logging.MaxSizeMB <= 0 {
		return errors.New("logging.max_size_mb must be greater than zero")
	}
	if cfg.Logging.MaxBackups < 0 {
		return errors.New("logging.max_backups must not be negative")
	}
	if cfg.Logging.MaxAgeDays < 0 {
		return errors.New("logging.max_age_days must not be negative")
	}
	if filepath.Clean(cfg.Storage.Path) == filepath.Clean(cfg.Logging.Path) {
		return errors.New("database and log paths must be different")
	}

	return nil
}
