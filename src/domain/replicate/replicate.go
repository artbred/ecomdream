package replicate

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	apiToken string
	baseURL  = "https://api.replicate.com/v1/predictions"
)

type Request struct {
	Version string                 `json:"version"`
	Input   map[string]interface{} `json:"input"`
}

type Response struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Urls    struct {
		Get    string `json:"get"`
		Cancel string `json:"cancel"`
	} `json:"urls"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	Status      string    `json:"status"`
	Input       any       `json:"input"`
	Output      []string    `json:"output"`
	Error       string    `json:"error"`
	Logs        string    `json:"logs"`
	Metrics     struct {
		PredictTime float64 `json:"predict_time"`
	} `json:"metrics"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	}

	apiToken = os.Getenv("REPLICATE_API_TOKEN")
}
