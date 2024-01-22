package routes

import (
	authcontroller "PBD_backend_go/controller/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(route fiber.Router) {
	route.Post("/login", authcontroller.LoginController)
}
