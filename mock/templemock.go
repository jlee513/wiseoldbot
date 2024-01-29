package mock

import (
	"context"
	"net/http"
	"osrs-disc-bot/util"
	"time"
)

type TempleClientMock struct {
	client *http.Client
}

func NewTempleClient() *TempleClientMock {
	client := new(TempleClientMock)
	client.client = &http.Client{Timeout: 30 * time.Second}
	return client
}

func (t *TempleClientMock) AddMemberToTemple(ctx context.Context, addingMember string) {
	return
}

func (t *TempleClientMock) RemoveMemberFromTemple(ctx context.Context, removingMember string) {
	return
}

func (t *TempleClientMock) GetKCsFromTemple(ctx context.Context, bossIdForTemple string) (*util.HallOfFameInfo, []int) {
	return nil, nil
}
