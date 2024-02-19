package routes

import (
	controller "PBD_backend_go/controller/inventory"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func InventoryRoute(routes fiber.Router) {
	routes.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Inventory",
			Name:      "CanView",
		})
	}, controller.GetInventoryController)
	routes.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Inventory",
			Name:      "CanEdit",
		})
	}, controller.GetInventoryByIDController)
	routes.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Inventory",
			Name:      "CanEdit",
		})
	}, controller.AddInventoryController)
	routes.Patch("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Inventory",
			Name:      "CanEdit",
		})
	}, controller.UpdateInventoryController)
	routes.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Inventory",
			Name:      "CanRemove",
		})
	}, controller.DeleteInventoryController)
}
