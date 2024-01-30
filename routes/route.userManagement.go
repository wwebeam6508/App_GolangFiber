package routes

import (
	authcontroller "PBD_backend_go/controller/userManagement"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func UserManagementRoute(route fiber.Router) {
	route.Get("/getUser", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "User",
			Name:      "CanView",
		})
	}, authcontroller.GetUserController)
	route.Get("/getUserByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "User",
			Name:      "CanEdit",
		})
	}, authcontroller.GetUserByIDController)
	route.Post("/addUser", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "User",
			Name:      "CanEdit",
		})
	}, middleware.RankCheck, authcontroller.AddUserController)
	route.Post("/updateUser", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "User",
			Name:      "CanEdit",
		})
	}, middleware.RankCheck, authcontroller.UpdateUserController)
	route.Delete("/deleteUser", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "User",
			Name:      "CanRemove",
		})
	}, middleware.RankCheck, authcontroller.DeleteUserController)
}
