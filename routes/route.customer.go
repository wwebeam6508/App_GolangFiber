package routes

import (
	controller "PBD_backend_go/controller/customer"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func CustomerRoute(route fiber.Router) {
	route.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanView",
		})
	}, controller.GetCustomerController)
	route.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanEdit",
		})
	}, controller.GetCustomerByIDController)
	route.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanEdit",
		})
	}, controller.AddCustomerController)
	route.Post("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanEdit",
		})
	}, controller.UpdateCustomerController)
	route.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanRemove",
		})
	}, controller.DeleteCustomerController)
}
