package service

import (
	"context"
	"io"
	"osrs-disc-bot/util"
)

// Holds the interfaces for all functions that make request externally

type imgur interface {
	Upload(ctx context.Context, AccessToken string, image io.Reader) string
	GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) string
}

type collectionLog interface {
	RetrieveCollectionLogAndOrder(ctx context.Context, submissions map[string]int) (map[string]int, []string)
}

type sheets interface {
	InitializeSubmissionsFromSheet(ctx context.Context, submissions map[string]int)
	UpdateCpSheet(ctx context.Context, submissions map[string]int)
	UpdateCpScreenshotsSheet(ctx context.Context, cpscreenshots map[string]string)
}

type temple interface {
	AddMemberToTemple(ctx context.Context, addingMember string, templeGroupId string, templeGroupKey string)
	RemoveMemberFromTemple(ctx context.Context, removingMember string, templeGroupId string, templeGroupKey string)
	GetPodiumFromTemple(ctx context.Context, bossIdForTemple string) (*util.HallOfFameInfo, []int)
}

type runescape interface {
	GetLeaguesPodiumFromRS(ctx context.Context, submissions map[string]int) (map[string]int, []string)
}
