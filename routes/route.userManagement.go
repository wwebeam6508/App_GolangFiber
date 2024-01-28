package routes

import (
	authcontroller "PBD_backend_go/controller/userManagement"
	"PBD_backend_go/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserManagementRoute(route fiber.Router) {
	route.Get("/getUser", middleware.Authenication, authcontroller.GetUserController)
	route.Get("/getUserByID", middleware.Authenication, authcontroller.GetUserByIDController)
	
}
