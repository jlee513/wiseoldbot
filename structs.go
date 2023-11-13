package main

type Config struct {
	DiscBotToken         string `mapstructure:"DISCORD_BOT_TOKEN"`
	DiscSubChan          string `mapstructure:"DISCORD_SUBMISSION_CHANNEL"`
	DiscSignUpChan       string `mapstructure:"DISCORD_SIGNUP_CHANNEL"`
	DiscLeaderboardChan  string `mapstructure:"DISCORD_LEADERBOARD_CHANNEL"`
	DiscDirectoryChan    string `mapstructure:"DISCORD_DIRECTORY_CHANNEL"`
	DiscSlayerBossesChan string `mapstructure:"DISCORD_SLAYER_BOSSES_CHANNEL"`
	DiscGwdChan          string `mapstructure:"DISCORD_GWD_CHANNEL"`
	DiscWildyChan        string `mapstructure:"DISCORD_WILDERNESS_CHANNEL"`
	DiscOtherChan        string `mapstructure:"DISCORD_OTHER_CHANNEL"`
	DiscMiscChan         string `mapstructure:"DISCORD_MISC_CHANNEL"`
	DiscDT2Chan          string `mapstructure:"DISCORD_DT2_CHANNEL"`
	DiscRaidsChan        string `mapstructure:"DISCORD_RAIDS_CHANNEL"`
	DiscPVPChan          string `mapstructure:"DISCORD_PVP_CHANNEL"`
	DiscCluesChan        string `mapstructure:"DISCORD_CLUES_CHANNEL"`
	DiscColChan          string `mapstructure:"DISCORD_COL_CHANNEL"`
	TempleGroupId        string `mapstructure:"TEMPLE_GROUP_ID"`
	TempleGroupKey       string `mapstructure:"TEMPLE_GROUP_KEY"`
	SheetsCp             string `mapstructure:"SHEETS_CP"`
	SheetsCpSub          string `mapstructure:"SHEETS_CP_SUB"`
	SheetsSpeedSub       string `mapstructure:"SHEETS_SPEED_SUB"`
}

type hallOfFameInfo struct {
	Data struct {
		BossName string         `json:"skill"`
		Players  map[int]Player `json:"players"`
	} `json:"data"`
}

type Player struct {
	Username string `json:"username"`
	Kc       int    `json:"xp"`
}

type hallOfFameRequestInfo struct {
	Bosses   map[string]string
	DiscChan string
}

type collectionLogInfo struct {
	CollectionLog struct {
		Uniques int `json:"uniqueObtained"`
	} `json:"collectionLog"`
}
