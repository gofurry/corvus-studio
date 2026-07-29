package cli

import (
	"context"
	"fmt"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/config"
	"github.com/spf13/cobra"
)

type RuntimeFunc func(context.Context, config.Config) error
type ConfigLoaderFunc func(config.LoadOptions) (config.Config, error)

type Dependencies struct {
	LoadConfig ConfigLoaderFunc
	Run        RuntimeFunc
}

type serveFlags struct {
	configFile  string
	dataDir     string
	host        string
	port        int
	storagePath string
	logPath     string
	logLevel    string
}

func NewRootCommand(dependencies Dependencies) *cobra.Command {
	if dependencies.LoadConfig == nil {
		dependencies.LoadConfig = config.Load
	}
	if dependencies.Run == nil {
		dependencies.Run = func(context.Context, config.Config) error {
			return fmt.Errorf("core runtime dependency is not configured")
		}
	}

	flags := &serveFlags{}
	root := &cobra.Command{
		Use:           "corvus",
		Short:         "Corvus Studio Core",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}
	root.PersistentFlags().StringVar(&flags.configFile, "config", "", "path to a YAML or JSON configuration file")

	serve := &cobra.Command{
		Use:   "serve",
		Short: "Start the local Corvus Core runtime",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			options := config.LoadOptions{ConfigFile: flags.configFile}
			overrides := &options.Overrides
			setStringOverride(command, "data-dir", flags.dataDir, &overrides.DataDir)
			setStringOverride(command, "host", flags.host, &overrides.Host)
			setIntOverride(command, "port", flags.port, &overrides.Port)
			setStringOverride(command, "database", flags.storagePath, &overrides.Storage)
			setStringOverride(command, "log-file", flags.logPath, &overrides.LogPath)
			setStringOverride(command, "log-level", flags.logLevel, &overrides.LogLevel)

			cfg, err := dependencies.LoadConfig(options)
			if err != nil {
				return fmt.Errorf("load Core configuration: %w", err)
			}
			if err := dependencies.Run(command.Context(), cfg); err != nil {
				return fmt.Errorf("run Core: %w", err)
			}
			return nil
		},
	}
	serve.Flags().StringVar(&flags.dataDir, "data-dir", "", "directory for the database, logs, and migration backups")
	serve.Flags().StringVar(&flags.host, "host", "", "loopback host to listen on")
	serve.Flags().IntVar(&flags.port, "port", 0, "TCP port to listen on")
	serve.Flags().StringVar(&flags.storagePath, "database", "", "SQLite database path")
	serve.Flags().StringVar(&flags.logPath, "log-file", "", "rotating log file path")
	serve.Flags().StringVar(&flags.logLevel, "log-level", "", "log level: debug, info, warn, or error")
	root.AddCommand(serve)

	return root
}

func setStringOverride(command *cobra.Command, name, value string, target **string) {
	if command.Flags().Changed(name) {
		copyOfValue := value
		*target = &copyOfValue
	}
}

func setIntOverride(command *cobra.Command, name string, value int, target **int) {
	if command.Flags().Changed(name) {
		copyOfValue := value
		*target = &copyOfValue
	}
}
