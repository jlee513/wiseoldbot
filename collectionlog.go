package main

import (
	"encoding/json"
	"fmt"
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func updateCollectionLog(session *discordgo.Session) {
	collectionLogs := make(map[string]int)
	for player, _ := range submissions {
		url := "https://api.collectionlog.net/collectionlog/user/" + player
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Unmarshal into hallOfFameInfo struct
			var col collectionLogInfo
			err = json.Unmarshal(bodyBytes, &col)
			if err != nil {
				fmt.Println("Error parsing JSON: ", err)
			}

			collectionLogs[player] = col.CollectionLog.Uniques
		} else if resp.StatusCode == http.StatusNotFound {
			fmt.Println("Missing user: " + player)
			continue
		}
	}

	// Update the #cp-leaderboard
	keys := make([]string, 0, len(collectionLogs))
	for key := range collectionLogs {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return collectionLogs[keys[i]] > collectionLogs[keys[j]]
	})

	// Create the leaderboard message that will be sent
	placements := ""
	for placement, k := range keys {
		switch placement {
		case 0:
			placements = placements + ":one: "
			break
		case 1:
			placements = placements + ":two: "
			break
		case 2:
			placements = placements + ":three: "
			break
		case 3:
			placements = placements + ":four: "
			break
		case 4:
			placements = placements + ":five: "
			break
		case 5:
			placements = placements + ":six: "
			break
		case 6:
			placements = placements + ":seven: "
			break
		case 7:
			placements = placements + ":eight: "
			break
		case 8:
			placements = placements + ":nine: "
			break
		case 9:
			placements = placements + ":keycap_10: "
			break

		}

		placements = placements + k + " [" + strconv.Itoa(collectionLogs[k]) + "]\n"
	}

	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(config.DiscColChan, 10, "", "", "")
	if err != nil {
		return
	}

	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	err = session.ChannelMessagesBulkDelete(config.DiscColChan, messageIDs)
	if err != nil {
		return
	}

	// Send the collection log message
	_, err = session.ChannelMessageSendEmbed(config.DiscColChan, embed.NewEmbed().
		SetTitle("Collection Log Ranking").
		SetDescription(placements).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/otTd8Dg.png").MessageEmbed)
	if err != nil {
		return
	}

	// Send the instructions on how to get on the collection log hall of fame
	var msg string
	msg = msg + "1. Download the collection-log plugin\n"
	msg = msg + "2. Click the box to \"Allow collectionlog.net connections\"\n"
	msg = msg + "3. Click through the collection log (there will be a * next to the one you still need to click)\n"
	msg = msg + "4. Go to the collection log icon on the sidebar\n"
	msg = msg + "5. Click Account at the top and then upload collection log\n"
	_, err = session.ChannelMessageSendEmbed(config.DiscColChan, embed.NewEmbed().
		SetTitle("How To Get Onto The Collection Log HOF").
		SetDescription(msg).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/otTd8Dg.png").MessageEmbed)
	if err != nil {
		return
	}
}
