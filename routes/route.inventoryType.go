package routes

import (
	controller "PBD_backend_go/controller/inventoryType"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func InventoryTypeRoute(routes fiber.Router) {
	routes.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "InventoryType",
			Name:      "CanView",
		})
	}, controller.GetInventoryTypeController)
	routes.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "InventoryType",
			Name:      "CanEdit",
		})
	}, controller.GetInventoryTypeByIDController)
	routes.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "InventoryType",
			Name:      "CanEdit",
		})
	}, controller.AddInventoryTypeController)
	routes.Put("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "InventoryType",
			Name:      "CanEdit",
		})
	}, controller.UpdateInventoryTypeController)
	routes.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "InventoryType",
			Name:      "CanRemove",
		})
	}, controller.DeleteInventoryTypeController)
	
}
