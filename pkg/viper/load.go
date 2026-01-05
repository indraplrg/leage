package viper

import (
	"share-notes-app/configs"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() (*configs.Config ,error) {
	// load json config
	viper.SetConfigName("app.config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs/")
	viper.AddConfigPath("/configs")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Warn("Config file not found")
		} else {
			return nil, err
		}
	} else {
		logrus.WithField("config_file", viper.ConfigFileUsed()).Info("Config file loaded")
	}

	var config configs.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}	
	
	// load env config
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &config, nil
} 