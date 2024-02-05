package routes

import (
	controller "PBD_backend_go/controller/userType"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func UserTypeRoute(route fiber.Router) {
	route.Get("/getUserType", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "UserType",
			Name:      "CanView",
		})
	}, controller.GetUserTypeController)
	route.Get("/getUserTypeByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "UserType",
			Name:      "CanEdit",
		})
	}, middleware.RankCheck, controller.GetUserTypeByIDController)
	route.Post("/addUserType", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "UserType",
			Name:      "CanEdit",
		})
	}, middleware.RankCheck, controller.AddUserTypeController)
	route.Post("/updateUserType", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "UserType",
			Name:      "CanEdit",
		})
	}, middleware.RankCheck, controller.UpdateUserTypeController)
	route.Delete("/deleteUserType", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "UserType",
			Name:      "CanRemove",
		})
	}, middleware.RankCheck, controller.DeleteUserTypeController)
}
