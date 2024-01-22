package controller

import (
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/auth"
	auth "PBD_backend_go/service/auth"

	"github.com/gofiber/fiber/v2"
)

func LoginController(c *fiber.Ctx) error {
	var body model.LoginRequest
	err := c.BodyParser(&body)
	if err != nil {
		exception.PanicLogging(err)
	}

	// call the LoginService function
	result, err := auth.LoginService(body)
	if err != nil {
		exception.PanicLogging(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
