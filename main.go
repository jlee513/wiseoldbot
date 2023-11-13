package main

import (
	"osrs-disc-bot/http"
	"osrs-disc-bot/service"
)

func main() {
	// Initialize the clients that make external calls
	collectionLog := http.NewCollectionLogClient()
	sheets := http.NewGoogleSheetsClient()
	imgur := http.NewImgurClient()
	temple := http.NewTempleClient()

	// Initialize configuration file
	config := initializeConfig()

	// Create the discord bot service and initialize the IRC
	osrsDiscBotService := service.NewService(config, collectionLog, sheets, imgur, temple)
	osrsDiscBotService.StartDiscordIRC()
}
