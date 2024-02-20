package routes

import (
	controller "PBD_backend_go/controller/requisition"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func RequisitionRoute(route fiber.Router) {
	route.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Requisition",
			Name:      "CanView",
		})
	}, controller.GetRequisitionController)
	route.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Requisition",
			Name:      "CanView",
		})
	}, controller.GetRequisitionByIDController)
	route.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Requisition",
			Name:      "CanEdit",
		})
	}, controller.AddRequisitionController)
	route.Put("/updateStatus", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Requisition",
			Name:      "CanEdit",
		})
	}, controller.UpdateRequisitionStatusService)
	route.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Requisition",
			Name:      "CanRemove",
		})
	}, controller.DeleteRequisitionController)
}
