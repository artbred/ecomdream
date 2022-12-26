package informer

import (
	"bytes"
	"ecomdream/src/pkg/configs"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	SendTelegramMessageEndpoint string = "/telegram/send-message"
	CallEndpoint                string = "/call"
)

type (
	SendTelegramMessageRequest struct {
		Token   string `json:"chat_token"`
		Message string `json:"message"`
	}

	CallRequest struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}

	JSONResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Client struct {
		BaseURL string
	}

	Level string
)

var (
	informerTokens = map[Level]string{}
	baseURL        string
)

func SendTelegramMessage(message string, level Level) {
	if configs.Debug {
		return
	}

	token, ok := informerTokens[level]
	if !ok {
		logrus.Warningf("Token for level %s is not set", level)
		return
	}

	url := baseURL + SendTelegramMessageEndpoint

	req := SendTelegramMessageRequest{
		Token:   token,
		Message: message,
	}

	b, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		return
	}

	logrus.Error("informer unsuccessfully")
}

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
		return
	}

	informerTokens = map[Level]string{
		"payments": os.Getenv("INFORMER_PAYMENTS_TOKEN"),
		"internal": os.Getenv("INFORMER_INTERNAL_TOKEN"),
	}

	baseURL = os.Getenv("INFORMER_BASE_URL")
}
