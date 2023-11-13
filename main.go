package main

import (
	"errors"
	"time"

	"github.com/do4mother/belajargo/authorization"
	"github.com/do4mother/belajargo/controller"
	"github.com/do4mother/belajargo/middleware"
	"github.com/do4mother/belajargo/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorResponse,
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var requestData LoginRequest

		if err := c.BodyParser(&requestData); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		users := []User{
			{
				Username: "admin",
				Password: "admin",
				Role:     "admin",
			},
			{
				Username: "user",
				Password: "user",
				Role:     "user",
			},
		}

		// get username and password from request body
		var loginRequest LoginRequest

		if err := c.BodyParser(&loginRequest); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		// find user
		findUser := func(username string) (User, error) {
			for _, user := range users {
				if user.Username == username {
					return user, nil
				}
			}
			return User{}, errors.New("User not found")
		}

		user, err := findUser(loginRequest.Username)

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// check username and password
		if user.Username != requestData.Username || user.Password != requestData.Password {
			return fiber.NewError(fiber.StatusUnauthorized, "Username or password is incorrect")
		}

		// generate jwt token
		claims := jwt.MapClaims{
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte("secret"))

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"token": t,
		})
	})

	proctected := app.Group("/protected", middleware.AuthProtected())
	proctected.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	proctected.Get("/me", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		return c.JSON(fiber.Map{
			"username": claims["username"],
			"role":     claims["role"],
		})
	})

	// testing authorization
	test := controller.NewTestAuthorizationController()
	proctected.Get("/test", test.Authorization.Middleware(authorization.Read, nil, nil), test.Get())
	proctected.Put("/test", test.Authorization.Middleware(authorization.Update, func(c *fiber.Ctx) (authorization.Map, error) {
		var testBody controller.TestBody

		if err := c.BodyParser(&testBody); err != nil {
			return nil, err
		}

		return authorization.Map{
			"username": testBody.Username,
		}, nil
	}, func() authorization.Map {
		// do query to database to compare with request body

		return authorization.Map{
			"username": "user",
		}
	}), test.Put())
	proctected.Post("/test", test.Authorization.Middleware(authorization.Create, nil, nil), test.Get())

	app.Listen(":8080")
}
