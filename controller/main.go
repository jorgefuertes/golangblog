package controller

import "github.com/gofiber/fiber/v2"

func Setup(app *fiber.App) {
	index(app)
	docs(app)
}
