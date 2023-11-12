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

type hallOfFameInfo struct {
	Data struct {
		BossName string         `json:"skill"`
		Players  map[int]Player `json:"players"`
	} `json:"data"`
}

type Player struct {
	Username string `json:"username"`
	Kc       int    `json:"xp"`
}

type hallOfFameRequestInfo struct {
	Bosses   map[string]string
	DiscChan string
}

func updateHallOfFame(session *discordgo.Session, requestInfo hallOfFameRequestInfo) {
	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", "", "")
	if err != nil {
		return
	}

	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
	if err != nil {
		return
	}

	// Now add all the bosses
	for bossIdForTemple, imageURL := range requestInfo.Bosses {
		getPodiumFromTemple(session, bossIdForTemple, requestInfo.DiscChan, imageURL)
	}
	return
}

func getPodiumFromTemple(session *discordgo.Session, bossIdForTemple string, discordChan string, imageURL string) {
	url := "https://templeosrs.com/api/skill_hiscores.php?group=2291&count=3&skill=" + bossIdForTemple

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
		var f hallOfFameInfo
		err = json.Unmarshal(bodyBytes, &f)
		if err != nil {
			fmt.Println("Error parsing JSON: ", err)
		}

		// Sort the map based on the keys
		keys := make([]int, 0, len(f.Data.Players))
		for key := range f.Data.Players {
			keys = append(keys, key)
		}
		sort.Ints(keys)

		// Iterate over the players to get the different places for users to create the placements
		placements := ""
		for _, k := range keys {
			switch k {
			case 1:
				placements = placements + ":first_place: "
				break
			case 2:
				placements = placements + ":second_place: "
				break
			case 3:
				placements = placements + ":third_place: "
				break
			}
			placements = placements + f.Data.Players[k].Username + " [" + strconv.Itoa(f.Data.Players[k].Kc) + "]\n"
		}
		
		_, err = session.ChannelMessageSendEmbed(discordChan, embed.NewEmbed().
			SetTitle(f.Data.BossName).
			SetDescription(placements).
			SetColor(0x1c1c1c).SetThumbnail(imageURL).MessageEmbed)
		if err != nil {
			return
		}
	}

}
