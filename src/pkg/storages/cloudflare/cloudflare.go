package cloudflare

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	api *cloudflare.API
	accountID, imagesBaseURL, email, token string
)

func init() {
	var err error

	err = godotenv.Load(); if err != nil {
		logrus.Error(err)
		return
	}

	accountID = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	email = os.Getenv("CLOUDFLARE_API_EMAIL")
	token = os.Getenv("CLOUDFLARE_API_TOKEN")

	imagesBaseURL = fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/images/v1", accountID)

	api, err = cloudflare.New(token, email)
	if err != nil {
		logrus.Error(err)
		return
	}
}
