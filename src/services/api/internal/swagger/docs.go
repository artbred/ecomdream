package swagger

import (
	_ "ecomdream/src/services/api/docs"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"os"
)

func ServeDocs(app *fiber.App) {
	app.Group("/swagger").
	Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			os.Getenv("ADMIN_USERNAME"): os.Getenv("ADMIN_PASSWORD"),
		},
	})).
	Get("*", swagger.HandlerDefault)
}
