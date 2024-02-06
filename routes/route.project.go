package routes

import (
	controller "PBD_backend_go/controller/project"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func ProjectRoute(route fiber.Router) {
	route.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Project",
			Name:      "CanView",
		})
	}, controller.GetProjectController)
	route.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Project",
			Name:      "CanEdit",
		})
	}, controller.GetProjectByIDController)
	route.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Project",
			Name:      "CanEdit",
		})
	}, controller.AddProjectController)
	route.Post("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Project",
			Name:      "CanEdit",
		})
	}, controller.UpdateProjectController)
	route.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Project",
			Name:      "CanRemove",
		})
	}, controller.DeleteProjectController)
	route.Get("/getCustomerName", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Customer",
			Name:      "CanView",
		})
	}, controller.GetCustomerNameController)

}
