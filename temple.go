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
	"strings"
)

func addNewMemberToTemple(newMember string) {
	url := "https://templeosrs.com/api/add_group_member.php"
	method := "POST"

	payload := strings.NewReader("id=" + config.TempleGroupId + "&key=" + config.TempleGroupKey + "&players=" + newMember)

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
func removeNewMemberToTemple(newMember string) {
	url := "https://templeosrs.com/api/remove_group_member.php"
	method := "POST"

	payload := strings.NewReader("id=" + config.TempleGroupId + "&key=" + config.TempleGroupKey + "&players=" + newMember)

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
