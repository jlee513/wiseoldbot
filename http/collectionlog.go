package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"osrs-disc-bot/util"
	"sort"
	"time"
)

type CollectionLogClient struct {
	client *http.Client
}

func NewCollectionLogClient() *CollectionLogClient {
	client := new(CollectionLogClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	return client
}

func (c CollectionLogClient) RetrieveCollectionLogAndOrder(submissions map[string]int) (map[string]int, []string) {
	collectionLog := make(map[string]int)
	for player, _ := range submissions {
		url := "https://api.collectionlog.net/collectionlog/user/" + player
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return nil, nil
			}

			// Unmarshal into hallOfFameInfo struct
			var col util.CollectionLogInfo
			err = json.Unmarshal(bodyBytes, &col)
			if err != nil {
				fmt.Println("Error parsing JSON: ", err)
				return nil, nil
			}

			collectionLog[player] = col.CollectionLog.Uniques
		} else if resp.StatusCode == http.StatusNotFound {
			fmt.Println("Missing user: " + player)
			continue
		}
	}
	// Update the #cp-leaderboard
	keys := make([]string, 0, len(collectionLog))
	for key := range collectionLog {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return collectionLog[keys[i]] > collectionLog[keys[j]]
	})

	return collectionLog, keys
}
