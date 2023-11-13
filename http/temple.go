package http

import (
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
	client *http.Client
}

func NewTempleClient() *TempleClient {
	client := new(TempleClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	return client
}

func (t *TempleClient) AddMemberToTemple(newMember string, templeGroupId string, templeGroupKey string) {
	url := "https://templeosrs.com/api/add_group_member.php"
	method := "POST"

	payload := strings.NewReader("id=" + templeGroupId + "&key=" + templeGroupKey + "&players=" + newMember)

	req, err := http.NewRequest(method, url, payload)
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

func (t *TempleClient) RemoveMemberFromTemple(removingMember string, templeGroupId string, templeGroupKey string) {
	url := "https://templeosrs.com/api/remove_group_member.php"
	method := "POST"

	payload := strings.NewReader("id=" + templeGroupId + "&key=" + templeGroupKey + "&players=" + removingMember)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
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

func (t *TempleClient) GetPodiumFromTemple(bossIdForTemple string) (*util.HallOfFameInfo, []int) {
	url := "https://templeosrs.com/api/skill_hiscores.php?group=2291&count=3&skill=" + bossIdForTemple

	resp, err := t.client.Get(url)
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
