package http

import (
	"context"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gemalto/flume"
)

type RunescapeClient struct {
	client     *http.Client
	leaguesUrl string
}

func NewRunescapeClient() *RunescapeClient {
	client := new(RunescapeClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.leaguesUrl = "https://secure.runescape.com/m=hiscore_oldschool_seasonal/index_lite.ws?player="
	return client
}

func (r *RunescapeClient) GetLeaguesPodiumFromRS(ctx context.Context, submissions map[string]int) (map[string]int, []string) {
	logger := flume.FromContext(ctx)
	leaguePodium := make(map[string]int)

	for player, _ := range submissions {
		resp, err := r.client.Get(r.leaguesUrl + player)
		if err != nil {
			logger.Error("Failed to get Runescape leagues info for: " + player)
			return nil, nil
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Failed to read body: " + err.Error())
				return nil, nil
			}

			// Set the body bytes as a string to parse through
			rsLeaguesInfo := string(bodyBytes[:])

			for lineNumber, leagueInfo := range strings.Split(rsLeaguesInfo, "\n") {
				if lineNumber == 24 {
					leaguePoints, _ := strconv.Atoi(strings.Split(leagueInfo, ",")[1])
					leaguePodium[player] = leaguePoints
				}
			}
		} else if resp.StatusCode == http.StatusNotFound {
			logger.Error("Missing user for leagues: " + player)
			continue
		}

	}

	// Create a slice and use that to sort the collectionLog
	keys := make([]string, 0, len(leaguePodium))
	for key := range leaguePodium {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return leaguePodium[keys[i]] > leaguePodium[keys[j]]
	})

	return leaguePodium, keys
}
