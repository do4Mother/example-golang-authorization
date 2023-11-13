package controller

import (
	"github.com/do4mother/example-golang-authorization/authorization"
	"github.com/gofiber/fiber/v2"
)

type TestAuthorization struct {
	Authorization *authorization.Authorization
}

func NewTestAuthorizationController() *TestAuthorization {
	a := new(authorization.Authorization)

	// admin
	a.Add("admin", authorization.Read, "*")

	// user
	a.Add("user", authorization.Read, "*")
	a.Add("user", authorization.Create, "*")
	a.Add("user", authorization.Update, "username")

	return &TestAuthorization{
		Authorization: a,
	}
}

func (t *TestAuthorization) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello world",
		})
	}
}

type TestBody struct {
	Username string `json:"username"`
}

func (t *TestAuthorization) Put() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello world",
		})
	}
}
