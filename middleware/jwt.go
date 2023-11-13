package middleware

import (
	"github.com/do4mother/example-golang-authorization/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthProtected() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte("secret"),
		},
		ErrorHandler: utils.ErrorResponse,
	})
}
