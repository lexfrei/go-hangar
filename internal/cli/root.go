package cli

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lexfrei/go-hangar/pkg/hangar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	baseURL      string
	apiToken     string
	timeout      time.Duration
	outputFormat string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "hangar",
	Short: "CLI tool for interacting with PaperMC Hangar API",
	Long: `hangar is a command-line interface for the PaperMC Hangar API.
It allows you to search for plugins, get version information, and download plugins.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(ctx context.Context) error {
	rootCmd.SetContext(ctx)
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/hangar/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", hangar.DefaultBaseURL, "Hangar API base URL")
	rootCmd.PersistentFlags().StringVar(&apiToken, "token", "", "Hangar API token")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", hangar.DefaultTimeout, "HTTP client timeout")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	// Bind flags to viper
	_ = viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	_ = viper.BindPFlag("api_token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Warn("failed to get home directory", "error", err)
			return
		}

		// Search config in home directory with name ".hangar" (without extension)
		viper.AddConfigPath(home + "/.config/hangar")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment variables
	viper.SetEnvPrefix("HANGAR")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("using config file", "file", viper.ConfigFileUsed())
	}
}

// createClient creates a new Hangar client from configuration.
func createClient() *hangar.Client {
	return hangar.NewClient(hangar.Config{
		BaseURL: viper.GetString("base_url"),
		Token:   viper.GetString("api_token"),
		Timeout: viper.GetDuration("timeout"),
	})
}
