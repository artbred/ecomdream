package replicate

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	apiToken string
	baseURL  = "https://api.replicate.com/v1/predictions"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	}

	apiToken = os.Getenv("REPLICATE_API_TOKEN")
}
