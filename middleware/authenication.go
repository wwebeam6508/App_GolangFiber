package middleware

import (
	"PBD_backend_go/exception"
	service "PBD_backend_go/service/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Authenication(c *fiber.Ctx) error {
	//get Authorization from header
	split := strings.Split(c.Get("Authorization"), " ")
	if len(split) != 2 {
		return exception.ValidationError{Message: "invalid token"}
	}
	token := split[1]
	//call verify jwt
	_, err := service.VerifyJWT(token)
	if err != nil {
		return exception.ErrorHandler(c, exception.UnauthorizedError{Message: "token is invalid or expired"})
	}
	return c.Next()
}
