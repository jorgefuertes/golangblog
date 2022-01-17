package controller

import (
	"golangblog/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func docs(app *fiber.App) {
	app.Get("/docs/cookies", func(c *fiber.Ctx) error {
		return helper.New(c, "docs/cookies").Render(nil)
	})
	app.Get("/docs/data", func(c *fiber.Ctx) error {
		return helper.New(c, "docs/data").Render(nil)
	})
}
