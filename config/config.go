package config

import (
	"github.com/spf13/viper"
	"osrs-disc-bot/util"
)

func InitializeConfig() *util.Config {
	var config util.Config
	// Initialize the Viper configuration ingestion and unmarshal
	viper.SetConfigFile("config/.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read environment variable, exiting now...")
	}
	err = viper.Unmarshal(&config)

	return &config
}
