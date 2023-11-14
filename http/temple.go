package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"osrs-disc-bot/util"
	"sort"
	"strings"
	"time"
)

type TempleClient struct {
	client       *http.Client
	addApiURL    string
	removeApiURL string
	podiumApiURL string
}

func NewTempleClient() *TempleClient {
	client := new(TempleClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.addApiURL = "https://templeosrs.com/api/add_group_member.php"
	client.removeApiURL = "https://templeosrs.com/api/remove_group_member.php"
	client.podiumApiURL = "https://templeosrs.com/api/skill_hiscores.php?group=2291&count=3&skill="
	return client
}

// AddMemberToTemple will make a POST request to the temple page to add a user to the group
func (t *TempleClient) AddMemberToTemple(ctx context.Context, addingMember string, templeGroupId string, templeGroupKey string) {
	payload := strings.NewReader("id=" + templeGroupId + "&key=" + templeGroupKey + "&players=" + addingMember)
	req, err := http.NewRequest(http.MethodPost, t.addApiURL, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = t.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// RemoveMemberFromTemple will make a POST request to the temple page to remove a user from the group
func (t *TempleClient) RemoveMemberFromTemple(ctx context.Context, removingMember string, templeGroupId string, templeGroupKey string) {
	payload := strings.NewReader("id=" + templeGroupId + "&key=" + templeGroupKey + "&players=" + removingMember)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, t.removeApiURL, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

/*
GetPodiumFromTemple will take in the bossid and make a request to temple to get the top 3 players
from our group with the highest kc
*/
func (t *TempleClient) GetPodiumFromTemple(ctx context.Context, bossIdForTemple string) (*util.HallOfFameInfo, []int) {
	resp, err := t.client.Get(t.podiumApiURL + bossIdForTemple)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil
		}

		// Unmarshal into hallOfFameInfo struct
		var f util.HallOfFameInfo
		err = json.Unmarshal(bodyBytes, &f)
		if err != nil {
			fmt.Println("Error parsing JSON: ", err)
			return nil, nil
		}

		// Sort the map based on the keys
		keys := make([]int, 0, len(f.Data.Players))
		for key := range f.Data.Players {
			keys = append(keys, key)
		}
		sort.Ints(keys)

		return &f, keys
	}

	return nil, nil
}
