package controller

import (
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/dashboard"
	service "PBD_backend_go/service/dashboard"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func TestDashboard(c *fiber.Ctx) error {
	var query model.GetEarnAndSpendEachYearInput
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}

	validate := validator.New()
	err := validate.Struct(query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	res, err := service.GetSpentAndEarnEachYear(query.Year)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.JSON(res)
}
