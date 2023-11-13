package main

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"os"
	"strconv"
)

func prepGoogleSheet(sheetId string) *spreadsheet.Sheet {
	// Create the client with the correct JWT configuration
	data, err := os.ReadFile("client_secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	// Fetch the clan points google sheet
	service := spreadsheet.NewServiceWithClient(client)
	googlesheet, err := service.FetchSpreadsheet(config.SheetsCp)
	checkError(err)
	sheet, err := googlesheet.SheetByIndex(0)
	checkError(err)

	return sheet
}

func initCpSheet() {
	sheet := prepGoogleSheet(config.SheetsCp)

	// Set the in memory submissions map with the Google sheets information
	submissions = make(map[string]int)
	for _, row := range sheet.Rows {
		isPlayer := true
		player := ""
		cp := 0
		for _, cell := range row {
			if isPlayer {
				player = cell.Value
			} else {
				cp, _ = strconv.Atoi(cell.Value)
				break
			}
			isPlayer = false
		}
		submissions[player] = cp
	}
}

func updateCpSheet() {
	sheet := prepGoogleSheet(config.SheetsCp)

	// Delete all the values in the sheet before proceeding with the insertion of clan points
	// We are deleting as this is an easier way of ensuring deleted people are removed from the
	// sheets without adding additional logic
	err := sheet.DeleteRows(0, len(sheet.Rows))
	checkError(err)

	// Update the Google sheets information with the in memory submissions map
	row := 0
	for player, cp := range submissions {
		sheet.Update(row, 0, player)
		sheet.Update(row, 1, strconv.Itoa(cp))
		row++
	}

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
