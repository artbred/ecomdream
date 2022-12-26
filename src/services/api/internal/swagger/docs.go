package swagger

import (
	_ "ecomdream/src/services/api/docs"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func ServeDocs(app *fiber.App) {
	app.Group("/swagger").Get("*", swagger.HandlerDefault)
}
