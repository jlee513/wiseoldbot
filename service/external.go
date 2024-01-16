package service

import (
	"context"
	"io"
	"osrs-disc-bot/util"
)

// Holds the interfaces for all functions that make request externally

type imgur interface {
	Upload(ctx context.Context, AccessToken string, image io.Reader) string
	GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) (string, error)
}

type collectionLog interface {
	RetrieveCollectionLogAndOrder(ctx context.Context, cp map[string]int) (map[string]int, []string)
}

type sheets interface {
	InitializeCpFromSheet(ctx context.Context, cp map[string]int)
	InitializeSpeedsFromSheet(ctx context.Context, speed map[string]util.SpeedInfo)
	InitializeFeedbackFromSheet(ctx context.Context, feedback map[string]string)
	InitializeTIDFromSheet(ctx context.Context) int
	UpdateTIDFromSheet(ctx context.Context, tid int)
	UpdateCpSheet(ctx context.Context, cp map[string]int)
	UpdateCpScreenshotsSheet(ctx context.Context, cpscreenshots map[string]string)
	UpdateSpeedSheet(ctx context.Context, speed map[string]util.SpeedInfo)
	UpdateSpeedScreenshotsSheet(ctx context.Context, speedscreenshots map[string]util.SpeedScInfo)
	UpdateFeedbackChannel(ctx context.Context, feedback map[string]string)
}

type temple interface {
	AddMemberToTemple(ctx context.Context, addingMember string, templeGroupId string, templeGroupKey string)
	RemoveMemberFromTemple(ctx context.Context, removingMember string, templeGroupId string, templeGroupKey string)
	GetPodiumFromTemple(ctx context.Context, bossIdForTemple string) (*util.HallOfFameInfo, []int)
}

type runescape interface {
	GetLeaguesPodiumFromRS(ctx context.Context, cp map[string]int) (map[string]int, []string)
}
