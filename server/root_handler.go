package server

import (
	"github.com/gofiber/fiber"
)

func RootHandler(c *fiber.Ctx) {
	c.JSON(
		map[string]interface{}{
			"message": "Method Not Allowed",
		},
	)
}

func RobotsTXTHandler(c *fiber.Ctx) {
	c.Send("User-agent: *\nDisallow: /")
}
