package controller

import (
	"github.com/do4mother/belajargo/authorization"
	"github.com/gofiber/fiber/v2"
)

type TestController struct {
	Authorization *authorization.Authorization
}

func NewTestController() *TestController {
	a := new(authorization.Authorization)

	// admin
	a.Add("admin", authorization.Read, "*")

	// user
	a.Add("user", authorization.Read, "*")
	a.Add("user", authorization.Create, "*")
	a.Add("user", authorization.Update, "username")

	return &TestController{
		Authorization: a,
	}
}

func (t *TestController) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello world",
		})
	}
}

type TestBody struct {
	Username string `json:"username"`
}

func (t *TestController) Put() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello world",
		})
	}
}
