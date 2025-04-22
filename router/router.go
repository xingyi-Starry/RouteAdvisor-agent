package router

import (
	"route-advisor-agent/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouter(app *fiber.App, c *services.Controller) {

	app.Group("/api/agent").
		Use(cors.New()).
		Post("/diary_writer", c.DiaryWriter.DiaryWriterHandler)
}
