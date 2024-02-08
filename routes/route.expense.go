package routes

import (
	controller "PBD_backend_go/controller/expense"
	"PBD_backend_go/middleware"
	"PBD_backend_go/model"

	"github.com/gofiber/fiber/v2"
)

func ExpenseRoute(route fiber.Router) {
	route.Get("/get", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Expense",
			Name:      "CanView",
		})
	}, controller.GetExpenseController)
	route.Get("/getByID", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Expense",
			Name:      "CanEdit",
		})
	}, controller.GetExpenseByIDController)
	route.Post("/add", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Expense",
			Name:      "CanEdit",
		})
	}, controller.AddExpenseController)
	route.Post("/update", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Expense",
			Name:      "CanEdit",
		})
	}, controller.UpdateExpenseController)
	route.Delete("/delete", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Expense",
			Name:      "CanRemove",
		})
	}, controller.DeleteExpenseController)
	route.Get("/getProjectTitle", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Expense",
			Name:      "CanEdit",
		})
	}, controller.GetWorkTitleController)
	route.Get("/getSellerName", middleware.Authenication, func(c *fiber.Ctx) error {
		return middleware.Permission(c, model.PermissionInput{
			GroupName: "Expense",
			Name:      "CanEdit",
		})
	}, controller.GetCustomerNameController)

}
