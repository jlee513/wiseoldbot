package service

import (
	"context"
	"io"
	"osrs-disc-bot/util"
)

// Holds the interfaces for all functions that make request externally

type imageservice interface {
	Upload(ctx context.Context, AccessToken string, image io.Reader) string
	GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) (string, error)
}

type collectionLog interface {
	RetrieveCollectionLogAndOrder(ctx context.Context, members map[string]util.MemberInfo) map[string]util.CollectionLogInfo
}

type sheets interface {
	InitializeCpFromSheet(ctx context.Context, cp map[string]int)
	InitializeSpeedsFromSheet(ctx context.Context, speed map[string]util.SpeedInfo)
	InitializeTIDFromSheet(ctx context.Context) int
	InitializeMembersFromSheet(ctx context.Context, members map[string]util.MemberInfo)
	UpdateMembersSheet(ctx context.Context, members map[string]util.MemberInfo)
	UpdateTIDFromSheet(ctx context.Context, tid int)
	UpdateCpSheet(ctx context.Context, cp map[string]int)
	UpdateCpScreenshotsSheet(ctx context.Context, cpscreenshots map[string]string)
	UpdateSpeedSheet(ctx context.Context, speed map[string]util.SpeedInfo)
	UpdateSpeedScreenshotsSheet(ctx context.Context, speedscreenshots map[string]util.SpeedScInfo)
}

type temple interface {
	AddMemberToTemple(ctx context.Context, addingMember string)
	RemoveMemberFromTemple(ctx context.Context, removingMember string)
	GetKCsFromTemple(ctx context.Context) *util.HallOfFameInfo
	GetMilestonesFromTemple(ctx context.Context) *util.MilestoneInfo
}

type runescape interface {
	GetLeaguesPodiumFromRS(ctx context.Context, cp map[string]int) (map[string]int, []string)
}

type pastebin interface {
	UpdateGuideList(ctx context.Context, pastebinGuides map[string][]util.GuideInfo)
	GetGuide(ctx context.Context, guideName string) string
}
