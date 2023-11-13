package service

import (
	"golang.org/x/net/context"
	"io"
	"osrs-disc-bot/util"
)

// Holds the interfaces for all functions that make request externally

type imgur interface {
	Upload(ctx context.Context, AccessToken string, image io.Reader) string
	GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) string
}

type collectionLog interface {
	RetrieveCollectionLogAndOrder(submissions map[string]int) (map[string]int, []string)
}

type sheets interface {
	InitializeSubmissionsFromSheet(submissions map[string]int)
	UpdateCpSheet(submissions map[string]int)
	UpdateCpScreenshotsSheet(cpscreenshots map[string]string)
}

type temple interface {
	AddMemberToTemple(addingMember string, templeGroupId string, templeGroupKey string)
	RemoveMemberFromTemple(removingMember string, templeGroupId string, templeGroupKey string)
	GetPodiumFromTemple(bossIdForTemple string) (*util.HallOfFameInfo, []int)
}
