package util

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
	DiscLeaguesChan      string `mapstructure:"DISCORD_LEAGUES_CHANNEL"`
	TempleGroupId        string `mapstructure:"TEMPLE_GROUP_ID"`
	TempleGroupKey       string `mapstructure:"TEMPLE_GROUP_KEY"`
	SheetsCp             string `mapstructure:"SHEETS_CP"`
	SheetsCpSC           string `mapstructure:"SHEETS_CP_SC"`
	SheetsSpeedSC        string `mapstructure:"SHEETS_SPEED_SC"`
	ImgurClientId        string `mapstructure:"IMGUR_CLIENT_ID"`
	ImgurClientSecret    string `mapstructure:"IMGUR_CLIENT_SECRET"`
	ImgurRefreshToken    string `mapstructure:"IMGUR_REFRESH_TOKEN"`
	LogDebug             bool   `mapstructure:"LOG_DEBUG"`
	CronKickoffTime      string `mapstructure:"CRON_KICKOFF_TIME"`
}

type PodiumLeaguePoints struct {
	Player string `json:"player"`
	Points int    `json:"points"`
}

type HallOfFameInfo struct {
	Data struct {
		BossName string         `json:"skill"`
		Players  map[int]Player `json:"players"`
	} `json:"data"`
}

type Player struct {
	Username string `json:"username"`
	Kc       int    `json:"xp"`
}

type HallOfFameRequestInfo struct {
	Name     string
	Bosses   map[string]string
	DiscChan string
}

type CollectionLogInfo struct {
	CollectionLog struct {
		Uniques int `json:"uniqueObtained"`
	} `json:"collectionLog"`
}

type ImageInfoDataWrapper struct {
	Ii      *ImageInfo `json:"data"`
	Success bool       `json:"success"`
}

// ImageInfo contains all image information provided by imgur
type ImageInfo struct {
	Link string `json:"link"` // The direct link to the image. (Note: if fetching an animated GIF that was over 20MB in original size, a .gif thumbnail will be returned)
}

type GenerateAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"` // The refresh token returned from the authorization code exchange
	ClientID     string `json:"client_id"`     // The client_id obtained during application registration
	ClientSecret string `json:"client_secret"` // The client secret obtained during application registration
	GrantType    string `json:"grant_type"`    // As defined in the OAuth2 specification, this field must contain a value of: refresh_token
}

type GenerateAccessTokenResponse struct {
	AccessToken     string `json:"access_token"` // TNew access token to use
	ExpiresIn       uint64 `json:"expires_in"`   // These parameters describe the lifetime of the token in seconds, and the kind of token that is being returned
	TokenType       string `json:"token_type"`
	Scope           string `json:"scope,omitempty"`            // Scope which were provided earlier during creation access_token
	RefreshToken    string `json:"refresh_token"`              // New refresh token
	AccountID       int    `json:"account_id,omitempty"`       // not specified in documentation
	AccountUserName string `json:"account_username,omitempty"` // not specified in documentation
}
