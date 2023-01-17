package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

var Debug bool

var StripeSecretKey string
var StripeWebhookSecret string

func init() {
	StripeSecretKey = os.Getenv("STRIPE_SECRET_KEY")
	StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
}
