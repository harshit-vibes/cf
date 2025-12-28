package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// AppName is the application name
	AppName = "dsaprep"

	// DefaultDailyGoal is the default daily problem goal
	DefaultDailyGoal = 5

	// DefaultMinRating is the default minimum difficulty
	DefaultMinRating = 800

	// DefaultMaxRating is the default maximum difficulty
	DefaultMaxRating = 1600
)

// Config holds the application configuration
type Config struct {
	CFHandle   string     `mapstructure:"cf_handle"`
	Difficulty Difficulty `mapstructure:"difficulty"`
	DailyGoal  int        `mapstructure:"daily_goal"`
	Theme      string     `mapstructure:"theme"`
}

// Difficulty holds difficulty range preferences
type Difficulty struct {
	Min int `mapstructure:"min"`
	Max int `mapstructure:"max"`
}

// configDir returns the configuration directory path
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "."+AppName), nil
}

// Init initializes the configuration
func Init(cfgFile string) error {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config directory
		configPath, err := configDir()
		if err != nil {
			return err
		}

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return err
		}

		// Search config in config directory
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Set defaults
	viper.SetDefault("cf_handle", "")
	viper.SetDefault("difficulty.min", DefaultMinRating)
	viper.SetDefault("difficulty.max", DefaultMaxRating)
	viper.SetDefault("daily_goal", DefaultDailyGoal)
	viper.SetDefault("theme", "dark")

	// Read environment variables
	viper.SetEnvPrefix("DSAPREP")
	viper.AutomaticEnv()

	// Read config file (ignore error if file doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

// Get returns the current configuration
func Get() *Config {
	return &Config{
		CFHandle: viper.GetString("cf_handle"),
		Difficulty: Difficulty{
			Min: viper.GetInt("difficulty.min"),
			Max: viper.GetInt("difficulty.max"),
		},
		DailyGoal: viper.GetInt("daily_goal"),
		Theme:     viper.GetString("theme"),
	}
}

// GetCFHandle returns the Codeforces handle
func GetCFHandle() string {
	return viper.GetString("cf_handle")
}

// GetDifficultyMin returns the minimum difficulty
func GetDifficultyMin() int {
	return viper.GetInt("difficulty.min")
}

// GetDifficultyMax returns the maximum difficulty
func GetDifficultyMax() int {
	return viper.GetInt("difficulty.max")
}

// GetDailyGoal returns the daily problem goal
func GetDailyGoal() int {
	return viper.GetInt("daily_goal")
}

// Set sets a configuration value
func Set(key string, value interface{}) {
	viper.Set(key, value)
}

// Save saves the current configuration to file
func Save() error {
	configPath, err := configDir()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configPath, "config.yaml")
	return viper.WriteConfigAs(configFile)
}

// GetConfigPath returns the config file path
func GetConfigPath() string {
	return viper.ConfigFileUsed()
}

// DataDir returns the data directory path
func DataDir() (string, error) {
	configPath, err := configDir()
	if err != nil {
		return "", err
	}
	return configPath, nil
}
