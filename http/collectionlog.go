package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"osrs-disc-bot/util"
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
cp map, sort it based on number of collection logs obtained, and return a map with the player's
name + collection log number along with the rankings
*/
func (c *CollectionLogClient) RetrieveCollectionLogAndOrder(ctx context.Context, members map[string]util.MemberInfo) map[string]util.CollectionLogInfo {
	logger := flume.FromContext(ctx)
	collectionLog := make(map[string]util.CollectionLogInfo)
	for player := range members {
		// Call the collectionlog api for the player
		url := "https://api.collectionlog.net/collectionlog/user/" + player
		resp, err := http.Get(url)
		if err != nil {
			logger.Error("Failed to contact collection log API: " + err.Error())
			return nil
		}
		defer resp.Body.Close()

		// Ensure we get an HTTP 200 response before unmarshaling
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Failed to read response body: " + err.Error())
				return nil
			}

			// Unmarshal into hallOfFameInfo struct
			var col util.CollectionLogInfo
			err = json.Unmarshal(bodyBytes, &col)
			if err != nil {
				logger.Error("Error parsing JSON: ", err.Error())
				return nil
			}

			// Remove all pets that have a quantity of zero
			var removedQuantityZero []util.PetInfo
			for _, petInfo := range col.CollectionLog.Tabs.Other.AllPets.Items {
				if petInfo.Quantity > 0 {
					removedQuantityZero = append(removedQuantityZero, petInfo)
				}
			}

			col.CollectionLog.Tabs.Other.AllPets.Items = removedQuantityZero

			// Set the collectionLog map with the player name as the key and the value as the
			// number of uniques this player has
			collectionLog[player] = col
		} else if resp.StatusCode == http.StatusNotFound {
			continue
		}
	}

	return collectionLog
}
