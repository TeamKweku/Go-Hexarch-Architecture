package config

import "github.com/spf13/viper"

type Config struct {
	DBDriver string `mapstructure:"CODE_ODESSEY_DB_DRIVER"`
	DBSource string `mapstructure:"CODE_ODESSEY_DATABASE_URI"`
}

func LoadConfig(path string) (config Config, err error) {
	// Add the directory where the .env file is located
	viper.AddConfigPath(path)

	// Set the config name to an empty string for .env files
	viper.SetConfigName(".env")

	// Set the config type to env
	viper.SetConfigType("env")

	// Automatically override values with environment variables
	viper.AutomaticEnv()

	// Read in the config file
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err

	}

	// Unmarshal the config into the Config struct
	err = viper.Unmarshal(&config)
	return config, err
}
