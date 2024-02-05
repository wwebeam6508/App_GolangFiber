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
}
