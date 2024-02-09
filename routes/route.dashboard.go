package routes

import (
	controller "PBD_backend_go/controller/dashboard"
	"PBD_backend_go/middleware"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoute(route fiber.Router) {
	route.Get("/testDashboard", middleware.Authenication, controller.TestDashboard)
}
