package routes

import (
	controller "PBD_backend_go/controller/location"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func LocationRoute(route fiber.Router) {
	route.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanView",
		})
	}, controller.GetLocationController)
	route.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanEdit",
		})
	}, controller.GetLocationByIDController)
	route.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanEdit",
		})
	}, controller.AddLocationController)
	route.Put("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanEdit",
		})
	}, controller.UpdateLocationController)
	route.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Location",
			Name:      "CanRemove",
		})
	}, controller.DeleteLocationController)
}
