package util

import "time"

type Config struct {
	LogDebug        bool   `mapstructure:"LOG_DEBUG"`
	CronKickoffTime string `mapstructure:"CRON_KICKOFF_TIME"`

	DiscGuildId           string `mapstructure:"DISCORD_GUILD_ID"`
	DiscBotToken          string `mapstructure:"DISCORD_BOT_TOKEN"`
	DiscSubChan           string `mapstructure:"DISCORD_SUBMISSION_CHANNEL"`
	DiscSignUpChan        string `mapstructure:"DISCORD_SIGNUP_CHANNEL"`
	DiscLeaderboardChan   string `mapstructure:"DISCORD_LEADERBOARD_CHANNEL"`
	DiscCpInformationChan string `mapstructure:"DISCORD_CP_INFORMATION_CHANNEL"`

	// Submission Approval Channels
	DiscCpApprovalChan    string `mapstructure:"DISCORD_CP_APPROVAL_CHANNEL"`
	DiscEventApprovalChan string `mapstructure:"DISCORD_EVENT_APPROVAL_CHANNEL"`
	DiscSpeedApprovalChan string `mapstructure:"DISCORD_SPEED_APPROVAL_CHANNEL"`

	// Third Party Information
	TempleGroupId     string `mapstructure:"TEMPLE_GROUP_ID"`
	TempleGroupKey    string `mapstructure:"TEMPLE_GROUP_KEY"`
	SheetsCp          string `mapstructure:"SHEETS_CP"`
	SheetsCpSC        string `mapstructure:"SHEETS_CP_SC"`
	SheetsSpeed       string `mapstructure:"SHEETS_SPEED"`
	SheetsSpeedSC     string `mapstructure:"SHEETS_SPEED_SC"`
	SheetsFeedback    string `mapstructure:"SHEETS_FEEDBACK"`
	ImgurClientId     string `mapstructure:"IMGUR_CLIENT_ID"`
	ImgurClientSecret string `mapstructure:"IMGUR_CLIENT_SECRET"`
	ImgurRefreshToken string `mapstructure:"IMGUR_REFRESH_TOKEN"`

	//Guide Update
	DiscGuideUpdateChan string `mapstructure:"DISCORD_GUIDE_UPDATE_CHANNEL"`

	// CM
	DiscTrioCMChinGuideChan  string `mapstructure:"DISCORD_TRIO_CM_CHIN_GUIDE_CHANNEL"`
	DiscTrioCMPrepGuideChan  string `mapstructure:"DISCORD_TRIO_CM_PREP_GUIDE_CHANNEL"`
	DiscTrioCMSurgeGuideChan string `mapstructure:"DISCORD_TRIO_CM_SURGE_GUIDE_CHANNEL"`
	DiscTrioCMUsefulInfoChan string `mapstructure:"DISCORD_TRIO_CM_USEFUL_INFO_CHANNEL"`

	// TOB
	DiscTobMageGuideChan  string `mapstructure:"DISCORD_TOB_MAGE_GUIDE_CHANNEL"`
	DiscTobRangeGuideChan string `mapstructure:"DISCORD_TOB_RANGE_GUIDE_CHANNEL"`
	DiscTobMeleeGuideChan string `mapstructure:"DISCORD_TOB_MELEE_GUIDE_CHANNEL"`

	// HOF
	DiscLootLogChan        string `mapstructure:"DISCORD_LOOT_LOG_CHANNEL"`
	DiscSlayerBossesChan   string `mapstructure:"DISCORD_SLAYER_BOSSES_CHANNEL"`
	DiscGwdChan            string `mapstructure:"DISCORD_GWD_CHANNEL"`
	DiscWildyChan          string `mapstructure:"DISCORD_WILDERNESS_CHANNEL"`
	DiscOtherChan          string `mapstructure:"DISCORD_OTHER_CHANNEL"`
	DiscMiscChan           string `mapstructure:"DISCORD_MISC_CHANNEL"`
	DiscDT2Chan            string `mapstructure:"DISCORD_DT2_CHANNEL"`
	DiscRaidsChan          string `mapstructure:"DISCORD_RAIDS_CHANNEL"`
	DiscPVPChan            string `mapstructure:"DISCORD_PVP_CHANNEL"`
	DiscCluesChan          string `mapstructure:"DISCORD_CLUES_CHANNEL"`
	DiscColChan            string `mapstructure:"DISCORD_COL_CHANNEL"`
	DiscLeaguesChan        string `mapstructure:"DISCORD_LEAGUES_CHANNEL"`
	DiscHOFLeaderboardChan string `mapstructure:"DISCORD_HOF_LEADERBOARD_CHANNEL"`

	// SPEED
	DiscSpeedTzhaarChan    string `mapstructure:"DISCORD_SPEED_TZHAAR_CHANNEL"`
	DiscSpeedSlayerChan    string `mapstructure:"DISCORD_SPEED_SLAYER_CHANNEL"`
	DiscSpeedNightmareChan string `mapstructure:"DISCORD_SPEED_NIGHTMARE_CHANNEL"`
	DiscSpeedNexChan       string `mapstructure:"DISCORD_SPEED_NEX_CHANNEL"`
	DiscSpeedSoloChan      string `mapstructure:"DISCORD_SPEED_SOLO_CHANNEL"`
	DiscSpeedCOXChan       string `mapstructure:"DISCORD_SPEED_COX_CHANNEL"`
	DiscSpeedTOBChan       string `mapstructure:"DISCORD_SPEED_TOB_CHANNEL"`
	DiscSpeedTOAChan       string `mapstructure:"DISCORD_SPEED_TOA_CHANNEL"`
	DiscSpeedAgilityChan   string `mapstructure:"DISCORD_SPEED_AGILITY_CHANNEL"`
}

type SpeedInfo struct {
	PlayersInvolved string    `json:"playersInvolved"`
	Time            time.Time `json:"time"`
	URL             string    `json:"url"`
}

type SpeedScInfo struct {
	PlayersInvolved string `json:"playersInvolved"`
	Time            string `json:"time"`
	BossName        string `json:"bossName"`
	URL             string `json:"submissionTime"`
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

type HofRequestInfo struct {
	Name     string
	Bosses   []BossInfo
	DiscChan string
	AfterId  string
}

type SpeedsRequestInfo struct {
	Name     string
	Bosses   []BossInfo
	DiscChan string
	AfterId  string
}

type BossInfo struct {
	BossName  string
	ImageLink string
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
