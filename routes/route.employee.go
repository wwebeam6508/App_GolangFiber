package routes

import (
	controller "PBD_backend_go/controller/employee"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func EmployeeRoute(route fiber.Router) {
	route.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Employee",
			Name:      "CanView",
		})
	}, controller.GetEmployeeController)
	route.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Employee",
			Name:      "CanEdit",
		})
	}, controller.GetEmployeeByIDController)
	route.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Employee",
			Name:      "CanEdit",
		})
	}, controller.AddEmployeeController)
}
