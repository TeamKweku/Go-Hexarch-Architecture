package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string `mapstructure:"CODE_ODESSEY_ENVIRONMENT"`
	AccessTokenDuration  string `mapstructure:"CODE_ODESSEY_TOKEN_DURATION"`
	RefreshTokenDuration string `mapstructure:"CODE_ODESSEY_REFRESH_TOKEN_DURATION"`
	DBDriver             string `mapstructure:"CODE_ODESSEY_DB_DRIVER"`
	DBSource             string `mapstructure:"CODE_ODESSEY_DATABASE_URI"`
	DBHost               string `mapstructure:"CODE_ODESSEY_DB_HOST"`
	DBPassword           string `mapstructure:"CODE_ODESSEY_DB_PASSWORD"`
	DBPort               string `mapstructure:"CODE_ODESSEY_DB_PORT"`
	DBName               string `mapstructure:"CODE_ODESSEY_DB_NAME"`
	DBSslMode            string `mapstructure:"CODE_ODESSEY_DB_SSL_MODE"`
	DBUser               string `mapstructure:"CODE_ODESSEY_DB_USER"`
	RPCPort              string `mapstructure:"CODE_ODESSEY_PORT"`

}

func LoadConfig(path string) (config Config, err error) {
	// Add the directory where the .env file is located
	viper.AddConfigPath(path)

	// For docker build
	viper.AddConfigPath("/app")

	viper.SetConfigName(".env")

	// Set the config type to env
	viper.SetConfigType("env")

	// Automatically override values with environment variables
	viper.AutomaticEnv()

	// Read in the config file
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the config into the Config struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return config, nil
}
