package main

import (
	"github.com/gemalto/flume"
	"github.com/spf13/viper"
)

type Config struct {
	DiscBotToken string `mapstructure:"DISCORD_BOT_TOKEN"`
	DiscServURL  string `mapstructure:"DISCORD_SERVER_URL"`
	DiscSubChan  string `mapstructure:"DISCORD_SUBMISSION_CHANNEL"`
}

var config Config

func init() {
	// Initialize the Viper configuration ingestion and unmarshal
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read environment variable, exiting now...")
		return
	}
	err = viper.Unmarshal(&config)
}

func main() {
	var log = flume.New("main")
	log.Info("TESTING")
	startDiscordIRC()
}
