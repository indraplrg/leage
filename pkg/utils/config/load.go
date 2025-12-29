package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig(logger *logrus.Logger) (*Config ,error) {
	// load json config
	viper.SetConfigName("app.config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs/")
	viper.AddConfigPath("../../configs/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn("Config file not found")
		} else {
			logger.WithError(err).Error("Error reading file config")
			return nil, fmt.Errorf("Error reading config file: %w", err)
		}
	} else {
		logger.WithField("config_file", viper.ConfigFileUsed()).Info("Config file loaded")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.WithError(err).Error("Failed unmarshal config")
		return nil, fmt.Errorf("Failed unmarsahl config: %w", err)
	}	
	
	// load env config
	if err := godotenv.Load(); err != nil {
		logger.WithError(err).Error("Failed to load .env file")
		return nil, fmt.Errorf("Failed to load .env file: %v", err)
	}

	logger.Info("Config loaded successfully")
	return &config, nil
} 