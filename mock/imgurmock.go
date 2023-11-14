package mock

import (
	"context"
	"io"
	"net/http"
	"time"
)

type ImgurClientMock struct {
	client *http.Client
}

func NewImgurClient() *ImgurClientMock {
	client := new(ImgurClientMock)
	client.client = &http.Client{Timeout: 30 * time.Second}
	return client
}

func (i *ImgurClientMock) Upload(ctx context.Context, AccessToken string, image io.Reader) string {
	return ""
}

func (i *ImgurClientMock) GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) string {
	return ""
}
