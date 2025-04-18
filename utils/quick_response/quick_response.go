package quick_response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ResponseTemplate struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func OK(c *fiber.Ctx) error {
	t := ResponseTemplate{
		Code:    200,
		Message: "OK",
	}
	return c.Status(200).JSON(t)
}

func BadRequest(c *fiber.Ctx, message string) error {
	t := ResponseTemplate{
		Code:    400,
		Message: message,
	}
	logrus.Error(message)
	return c.Status(400).JSON(t)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	t := ResponseTemplate{
		Code:    401,
		Message: message,
	}
	logrus.Error(message)
	return c.Status(401).JSON(t)
}

func NotFound(c *fiber.Ctx, message string) error {
	t := ResponseTemplate{
		Code:    404,
		Message: message,
	}
	logrus.Error(message)
	return c.Status(404).JSON(t)
}

func Internal(c *fiber.Ctx, message string) error {
	t := ResponseTemplate{
		Code:    500,
		Message: "Internal Server Error",
	}
	logrus.Error(message)
	return c.Status(500).JSON(t)
}
