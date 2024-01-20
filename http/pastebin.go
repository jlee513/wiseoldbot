package http

import (
	"context"
	"github.com/TwiN/go-pastebin"
	"net/http"
	"strings"
	"time"
)

type PastebinClient struct {
	client         *http.Client
	pastebinClient *pastebin.Client

	pastebinUsername     string
	pastebinPassword     string
	pastebinDevApiKey    string
	pastebinMainPasteKey string

	pastebinGuides map[string]string
}

func NewPastebinClient(pastebinUsername, pastebinPassword, pastebinDevApiKey, pastebinMainPasteKey string) *PastebinClient {
	client := new(PastebinClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.pastebinUsername = pastebinUsername
	client.pastebinPassword = pastebinPassword
	client.pastebinDevApiKey = pastebinDevApiKey
	client.pastebinMainPasteKey = pastebinMainPasteKey
	client.pastebinGuides = make(map[string]string)

	pastebinClient, err := pastebin.NewClient(pastebinUsername, pastebinPassword, pastebinDevApiKey)
	if err != nil {
		panic(err)
	}

	client.pastebinClient = pastebinClient

	// Initialize the list of guides
	pasteContent, err := pastebinClient.GetUserPasteContent(pastebinMainPasteKey)
	if err != nil {
		panic(err)
	}

	guides := strings.Split(pasteContent, "\n")
	for _, guide := range guides {
		guideNameAndKey := strings.Split(guide, " - ")
		client.pastebinGuides[guideNameAndKey[0]] = guideNameAndKey[1]
	}

	return client
}

func (p *PastebinClient) GetGuide(ctx context.Context, guideName string) string {
	pasteContent, err := p.pastebinClient.GetUserPasteContent(p.pastebinGuides[guideName])
	if err != nil {
		panic(err)
	}
	return pasteContent
}
