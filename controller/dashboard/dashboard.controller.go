package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/dashboard"
	service "PBD_backend_go/service/dashboard"

	"github.com/gofiber/fiber/v2"
)

func TestDashboard(c *fiber.Ctx) error {
	var query model.GetEarnAndSpendEachYearInput
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}

	err := common.Validate(query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	res, err := service.GetTotalEarn(query.Year)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.JSON(res)
}
