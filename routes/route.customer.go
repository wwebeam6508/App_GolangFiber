package routes

import (
	controller "PBD_backend_go/controller/customer"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func CustomerRoute(route fiber.Router) {
	route.Get("/getCustomer", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanView",
		})
	}, controller.GetCustomerController)
	route.Get("/getCustomerByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanEdit",
		})
	}, controller.GetCustomerByIDController)
	route.Post("/addCustomer", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanEdit",
		})
	}, controller.AddCustomerController)
	route.Post("/updateCustomer", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanEdit",
		})
	}, controller.UpdateCustomerController)
}
