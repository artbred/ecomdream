package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type ImageUploadRequestByURL struct {
	URL string
	RequireSignedURLs bool
	Metadata map[string]interface{}
}

func UploadImageByURL(imageFromURL ImageUploadRequestByURL) (imageResponse *cloudflare.ImageDetailsResponse, err error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	err = w.WriteField("url", imageFromURL.URL)
	err = w.WriteField("requireSignedURLs", strconv.FormatBool(imageFromURL.RequireSignedURLs))

	if imageFromURL.Metadata != nil {
		part, err := w.CreateFormField("metadata")
		if err != nil {
			return nil, fmt.Errorf("create from field metadata, %w", err)
		}

		err = json.NewEncoder(part).Encode(imageFromURL.Metadata)
		if err != nil {
			return nil, err
		}
	}

	err = w.Close(); if err != nil {
		return nil, fmt.Errorf("cloudflare write multipart, %w", err)
	}

	req, _ := http.NewRequest("POST", imagesBaseURL, body)
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("X-Auth-Key", token)
	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("Accept","application/json")

	res, err := http.DefaultClient.Do(req); if err != nil {
		return nil, fmt.Errorf("can't upload image to cloudflare, %w", err)
	}

	bodyRes, err := io.ReadAll(res.Body); if err != nil {
		return nil, fmt.Errorf("io.ReadALL cloudlfare response: %w", err)
	}

	imageResponse = &cloudflare.ImageDetailsResponse{}
	err = json.Unmarshal(bodyRes, imageResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshall cloudlfare response: %w", err)
	}

	if imageResponse.Success != true {
		return nil, fmt.Errorf("imageResponse not success for url %s, %+v", imageFromURL.URL, imageResponse.Errors)
	}

	return
}
