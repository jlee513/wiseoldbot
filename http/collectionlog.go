package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"osrs-disc-bot/util"
	"sort"
	"time"

	"github.com/gemalto/flume"
)

type CollectionLogClient struct {
	client *http.Client
}

func NewCollectionLogClient() *CollectionLogClient {
	client := new(CollectionLogClient)
	client.client = &http.Client{Timeout: 30 * time.Second}

	return client
}

/*
RetrieveCollectionLogAndOrder will call the collectionlog.net's api for each of the players in the
submissions map, sort it based on number of collection logs obtained, and return a map with the player's
name + collection log number along with the rankings
*/
func (c CollectionLogClient) RetrieveCollectionLogAndOrder(ctx context.Context, submissions map[string]int) (map[string]int, []string) {
	logger := flume.FromContext(ctx)
	collectionLog := make(map[string]int)
	for player, _ := range submissions {
		// Call the collectionlog api for the player
		url := "https://api.collectionlog.net/collectionlog/user/" + player
		resp, err := http.Get(url)
		if err != nil {
			logger.Error("Failed to contact collection log API: " + err.Error())
			return nil, nil
		}
		defer resp.Body.Close()

		// Ensure we get an HTTP 200 response before unmarshaling
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Failed to read response body: " + err.Error())
				return nil, nil
			}

			// Unmarshal into hallOfFameInfo struct
			var col util.CollectionLogInfo
			err = json.Unmarshal(bodyBytes, &col)
			if err != nil {
				logger.Error("Error parsing JSON: ", err.Error())
				return nil, nil
			}

			// Set the collectionLog map with the player name as the key and the value as the
			// number of uniques this player has
			collectionLog[player] = col.CollectionLog.Uniques
		} else if resp.StatusCode == http.StatusNotFound {
			logger.Error("Missing user for collection log: " + player)
			continue
		}
	}

	// Create a slice and use that to sort the collectionLog
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
