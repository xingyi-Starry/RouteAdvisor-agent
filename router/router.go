package router

import (
	"route-advisor-agent/services"

	"github.com/gofiber/fiber/v2"
)

func InitRouter(app *fiber.App, c *services.Controller) {

	app.Group("/api/agent").
		Post("/diary_writer", c.DiaryWriter.DiaryWriterHandler)
}
