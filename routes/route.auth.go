package routes

import (
	authcontroller "PBD_backend_go/controller/auth"
	"PBD_backend_go/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(route fiber.Router) {
	route.Post("/login", authcontroller.LoginController)
	route.Post("/refreshtoken", authcontroller.RefreshTokenController)
	route.Get("/fetchuser", middleware.Authenication, authcontroller.FetchUserController)
	route.Post("/logout", authcontroller.LogoutController)
	route.Post("changepassword", authcontroller.ChangePasswordController)
}
