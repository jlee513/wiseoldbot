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
	sheets := http.NewGoogleSheetsClient(cfg)
	collectionLog := http.NewCollectionLogClient()
	imgur := http.NewImgurClient()
	temple := http.NewTempleClient()
	runescape := http.NewRunescapeClient()
	pastebin := http.NewPastebinClient(cfg.PastebinUsername, cfg.PastebinPassword, cfg.PastebinDevApiKey, cfg.PastebinMainPasteKey)

	// Create the discord bot service and initialize the IRC
	osrsDiscBotService := service.NewService(cfg, collectionLog, sheets, imgur, temple, runescape, pastebin)
	osrsDiscBotService.StartDiscordIRC()
}
