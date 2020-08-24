package config

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

//AppConfig Application configuration
type AppConfig struct {
	Port     int `yaml:"port"`
	Database struct {
		Driver   string `yaml:"driver"`
		Name     string `yaml:"name"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

//InitConfig Initiatilize config
func InitConfig() *AppConfig {
	var defaultConfig AppConfig
	defaultConfig.Port = 1323
	defaultConfig.Database.Driver = "mongodb"
	defaultConfig.Database.Name = "transaction"
	defaultConfig.Database.Address = "localhost"
	defaultConfig.Database.Port = 27017
	defaultConfig.Database.Username = ""
	defaultConfig.Database.Password = ""

	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("config file not found, will use default value")
		} else {
			log.Info("error to load config file, will use default value")
		}

		return &defaultConfig
	}

	var finalConfig AppConfig
	err := viper.Unmarshal(&finalConfig)
	if err != nil {
		log.Info("failed to extract config, will use default value")
		return &defaultConfig
	}

	return &finalConfig
}
