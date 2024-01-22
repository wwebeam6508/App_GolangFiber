package controller

import (
	"PBD_backend_go/commonentity"
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
		return c.Status(fiber.ErrNotAcceptable.Code).JSON(commonentity.GeneralResponse{
			Code:    fiber.ErrNotAcceptable.Code,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})

}
