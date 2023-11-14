package mock

import (
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

func (t *TempleClientMock) AddMemberToTemple(addingMember string, templeGroupId string, templeGroupKey string) {
	return
}

func (t *TempleClientMock) RemoveMemberFromTemple(removingMember string, templeGroupId string, templeGroupKey string) {
	return
}

func (t *TempleClientMock) GetPodiumFromTemple(bossIdForTemple string) (*util.HallOfFameInfo, []int) {
	return nil, nil
}
