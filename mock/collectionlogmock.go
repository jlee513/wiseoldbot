package mock

import (
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"
)

type CollectionLogClientMock struct {
	client *http.Client
}

func NewCollectionLogClient() *CollectionLogClientMock {
	client := new(CollectionLogClientMock)
	client.client = &http.Client{Timeout: 30 * time.Second}
	return client
}

func (c *CollectionLogClientMock) Upload(ctx context.Context, AccessToken string, image io.Reader) string {
	return ""
}

func (c *CollectionLogClientMock) GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) string {
	return ""
}
