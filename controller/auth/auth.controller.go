package controller

import (
	"PBD_backend_go/exception"
	service "PBD_backend_go/service/auth"

	"github.com/gofiber/fiber/v2"
)

func LoginController(c *fiber.Ctx) error {
	var body LoginRequest
	err := c.BodyParser(&body)
	if err != nil {
		exception.PanicLogging(err)
	}

	// call the LoginService function
	result, err := service.LoginService(body.username, body.password)
	if err != nil {
		exception.PanicLogging(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
