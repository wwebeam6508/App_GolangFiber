package routes

import (
	controller "PBD_backend_go/controller/wage"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func WageRoute(route fiber.Router) {
	route.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Wage",
			Name:      "CanView",
		})
	}, controller.GetWageController)
	route.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Wage",
			Name:      "CanEdit",
		})
	}, controller.GetWageByIDController)
	route.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Wage",
			Name:      "CanEdit",
		})
	}, controller.AddWageController)
	route.Put("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Wage",
			Name:      "CanEdit",
		})
	}, controller.UpdateWageController)
	route.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Wage",
			Name:      "CanRemove",
		})
	}, controller.DeleteWageController)

}
