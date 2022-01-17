package controller

import (
	"golangblog/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func index(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return helper.New(c, "index").SetTitle("Últimas publicaciones").Render(nil)
	})
}
