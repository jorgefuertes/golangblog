package helper

import (
	"github.com/gofiber/fiber/v2"
)

type Helper struct {
	c         *fiber.Ctx
	Template  string
	Theme     string
	Title     string
	MyURL     string
	CookiesOn bool
	Locals    *fiber.Map
}

func New(c *fiber.Ctx, tpl string) *Helper {
	h := new(Helper)
	h.c = c
	h.Template = tpl
	h.Theme = "amber"
	h.MyURL = c.BaseURL()

	return h
}

func (h *Helper) SetTitle(title string) *Helper {
	h.Title = title
	return h
}

func (h *Helper) Render(locals interface{}) error {
	if locals == nil {
		locals = &fiber.Map{}
	}
	return h.c.Render(h.Template, fiber.Map{"h": h, "l": locals})
}
