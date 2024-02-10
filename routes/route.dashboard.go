package routes

import (
	controller "PBD_backend_go/controller/dashboard"
	"PBD_backend_go/middleware"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoute(route fiber.Router) {
	route.Get("/getDashboard", middleware.Authenication, controller.GetDashboardController)
	route.Get("/getEarnAndSpendEachYear", middleware.Authenication, controller.GetEarnAndSpendEachYearController)
	route.Get("/getTotalEarn", middleware.Authenication, controller.GetTotalEarnController)
	route.Get("/getTotalExpense", middleware.Authenication, controller.GetTotalExpenseController)
	route.Get("/getTotalWork", middleware.Authenication, controller.GetTotalWorkController)
	route.Get("/getYearsReport", middleware.Authenication, controller.GetYearsReportController)
}
