package routes

import (
	controller "PBD_backend_go/controller/location"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func LocationRoute(routes fiber.Router) {
	routes.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanView",
		})
	}, controller.GetLocationController)
	routes.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanEdit",
		})
	}, controller.GetLocationByIDController)
	routes.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanEdit",
		})
	}, controller.AddLocationController)
	routes.Patch("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanEdit",
		})
	}, controller.UpdateLocationController)
	routes.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanRemove",
		})
	}, controller.DeleteLocationController)
}
