package main

import (
	"osrs-disc-bot/config"
	"osrs-disc-bot/http"
	"osrs-disc-bot/service"
)

func main() {
	// Initialize configuration file
	cfg := config.InitializeConfig()

	// Initialize the clients that make external calls
	collectionLog := http.NewCollectionLogClient()
	sheets := http.NewGoogleSheetsClient(cfg.SheetsCp, cfg.SheetsCpSC, cfg.SheetsSpeed, cfg.SheetsSpeedSC)
	imgur := http.NewImgurClient()
	temple := http.NewTempleClient()
	runescape := http.NewRunescapeClient()

	// Create the discord bot service and initialize the IRC
	osrsDiscBotService := service.NewService(cfg, collectionLog, sheets, imgur, temple, runescape)
	osrsDiscBotService.StartDiscordIRC()
}
