package http

import (
	"context"
	"github.com/TwiN/go-pastebin"
	"github.com/gemalto/flume"
	"net/http"
	"osrs-disc-bot/util"
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

	pastebinGuides map[string][]util.GuideInfo
}

func NewPastebinClient(pastebinUsername, pastebinPassword, pastebinDevApiKey, pastebinMainPasteKey string) *PastebinClient {
	client := new(PastebinClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.pastebinUsername = pastebinUsername
	client.pastebinPassword = pastebinPassword
	client.pastebinDevApiKey = pastebinDevApiKey
	client.pastebinMainPasteKey = pastebinMainPasteKey
	client.pastebinGuides = make(map[string][]util.GuideInfo)

	//pastebinClient, err := pastebin.NewClient(pastebinUsername, pastebinPassword, pastebinDevApiKey)
	//if err != nil {
	//	panic(err)
	//}
	//
	//client.pastebinClient = pastebinClient
	//
	//// Initialize the list of guides
	//client.UpdateGuideList(context.Background(), client.pastebinGuides)

	return client
}

func (p *PastebinClient) UpdateGuideList(ctx context.Context, pastebinGuides map[string][]util.GuideInfo) {
	//pasteContent, err := p.pastebinClient.GetUserPasteContent(p.pastebinMainPasteKey)
	//if err != nil {
	//	flume.FromContext(ctx).Error("Failed to get user pastebin content: " + err.Error())
	//}
	//
	//// Ensure everything gets deleted when updating
	//for guide := range pastebinGuides {
	//	delete(pastebinGuides, guide)
	//}
	//
	//pasteContent = strings.Replace(pasteContent, "\r", "", -1)
	//guides := strings.Split(pasteContent, "\n")
	//currentGuideName := ""
	//var currentGuideInfos []util.GuideInfo
	//for _, guide := range guides {
	//	// This is a top level guide and anything under this until the next # are different channels as a part
	//	// of the guide
	//	if strings.Contains(guide, "# ") {
	//		if len(currentGuideInfos) > 0 {
	//			pastebinGuides[currentGuideName] = currentGuideInfos
	//		}
	//		currentGuideName = strings.Replace(guide, "# ", "", -1)
	//		currentGuideInfos = nil
	//	} else {
	//		guideNameAndKey := strings.Split(guide, " - ")
	//		currentGuideInfos = append(currentGuideInfos, util.GuideInfo{
	//			GuidePageName: guideNameAndKey[0],
	//			PastebinKey:   guideNameAndKey[1],
	//			DiscChan:      guideNameAndKey[2],
	//		})
	//	}
	//}
	//// Ensure we get the last guides
	//pastebinGuides[currentGuideName] = currentGuideInfos
}

func (p *PastebinClient) GetGuide(ctx context.Context, guideKey string) string {
	pasteContent, err := p.pastebinClient.GetUserPasteContent(guideKey)
	if err != nil {
		flume.FromContext(ctx).Error("Failed to get paste content: " + err.Error())
	}
	return strings.Replace(pasteContent, "\r", "", -1)
}
