package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"osrs-disc-bot/util"
	"strings"
	"time"

	"github.com/gemalto/flume"
)

type TempleClient struct {
	client          *http.Client
	addApiURL       string
	removeApiURL    string
	groupMemberInfo string
	milestoneApiURL string
	templeGroupId   string
	templeGroupKey  string
}

func NewTempleClient(templeGroupId, templeGroupKey string) *TempleClient {
	client := new(TempleClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.addApiURL = "https://templeosrs.com/api/add_group_member.php"
	client.removeApiURL = "https://templeosrs.com/api/remove_group_member.php"
	client.groupMemberInfo = "https://templeosrs.com/api/group_member_info.php?bosses=true&id=" + templeGroupId
	client.milestoneApiURL = "https://templeosrs.com/api/group_achievements.php?id=" + templeGroupId
	client.templeGroupId = templeGroupId
	client.templeGroupKey = templeGroupKey
	return client
}

// AddMemberToTemple will make a POST request to the temple page to add a user to the group
func (t *TempleClient) AddMemberToTemple(ctx context.Context, addingMember string) {
	logger := flume.FromContext(ctx)
	logger.Info("Attempting to add new user to temple group: " + addingMember)

	payload := strings.NewReader("id=" + t.templeGroupId + "&key=" + t.templeGroupKey + "&players=" + addingMember)
	req, err := http.NewRequest(http.MethodPost, t.addApiURL, payload)
	if err != nil {
		logger.Error("Error while creating request: ", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = t.client.Do(req)
	if err != nil {
		logger.Error("Error while executing call to temple API: ", err.Error())
		return
	}
}

// RemoveMemberFromTemple will make a POST request to the temple page to remove a user from the group
func (t *TempleClient) RemoveMemberFromTemple(ctx context.Context, removingMember string) {
	logger := flume.FromContext(ctx)
	logger.Info("Attempting to remove user from temple group: " + removingMember)

	payload := strings.NewReader("id=" + t.templeGroupId + "&key=" + t.templeGroupKey + "&players=" + removingMember)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, t.removeApiURL, payload)
	if err != nil {
		logger.Error("Error while creating request: ", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(req)
	if err != nil {
		logger.Error("Error while executing call to temple API: ", err.Error())
		return
	}
}

/*
GetKCsFromTemple will take in the bossid and make a request to temple to get the kc of all the players
from our group
*/
func (t *TempleClient) GetKCsFromTemple(ctx context.Context) *util.HallOfFameInfo {
	logger := flume.FromContext(ctx)

	resp, err := t.client.Get(t.groupMemberInfo)
	if err != nil {
		logger.Error("Error while retrieving stats from temple API: ", err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Error while reading response from temple: ", err.Error())
			return nil
		}

		// Unmarshal into hallOfFameInfo struct
		var f util.HallOfFameInfo
		err = json.Unmarshal(bodyBytes, &f)
		if err != nil {
			logger.Error("Error parsing JSON: ", err.Error())
			return nil
		}

		return &f
	}

	return nil
}

func (t *TempleClient) GetMilestonesFromTemple(ctx context.Context) *util.MilestoneInfo {
	logger := flume.FromContext(ctx)

	resp, err := t.client.Get(t.milestoneApiURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while retrieving milestones from temple API"), err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Error while reading response from temple: ", err.Error())
			return nil
		}

		// Unmarshal into hallOfFameInfo struct
		var f util.MilestoneInfo
		err = json.Unmarshal(bodyBytes, &f)
		if err != nil {
			logger.Error("Error parsing JSON: ", err.Error())
			return nil
		}

		return &f
	}

	return nil
}
