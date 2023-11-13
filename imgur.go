package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

func uploadToImgur(submissionLink string) string {
	// Download the submission link in order to get the bytes in order to upload to imgur
	client := &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Get the image in the response and send that to the imgurUpload function
	resp, err := client.Get(submissionLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	imageUrl := imgurUpload(resp.Body)
	return imageUrl
}

func imgurUpload(image io.Reader) string {
	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, _ := writer.CreateFormFile("image", "dont care about name")
	io.Copy(part, image)

	writer.Close()
	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	token := getNewAccessToken()
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()

	// Read the whole body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	var img imageInfoDataWrapper
	if err = dec.Decode(&img); err != nil {
		panic(err)
	}

	return img.Ii.Link
}

func getNewAccessToken() string {
	client := &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	rawBody, err := json.Marshal(
		GenerateAccessTokenRequest{
			RefreshToken: config.ImgurRefreshToken,
			ClientID:     config.ImgurClientId,
			ClientSecret: config.ImgurClientSecret,
			GrantType:    "refresh_token",
		})

	req, err := http.NewRequest(http.MethodPost, "https://api.imgur.com/oauth2/token", bytes.NewBuffer(rawBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response := GenerateAccessTokenResponse{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err = decoder.Decode(&response); err != nil {
		panic(err)
	}

	return response.AccessToken
}
