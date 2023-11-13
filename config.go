package main

import (
	"github.com/spf13/viper"
	"osrs-disc-bot/util"
)

var config util.Config

func initializeConfig() *util.Config {
	// Initialize the Viper configuration ingestion and unmarshal
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read environment variable, exiting now...")
	}
	err = viper.Unmarshal(&config)

	return &config
}
