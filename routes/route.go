package routes

import (

	// import other routes...

	"PBD_backend_go/commonentity"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// print all error line of code
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stackTrace := debug.Stack()
				log.Printf("Panic: %v\nTrace:\n%s", err, stackTrace)

				c.Status(500).JSON(commonentity.GeneralResponse{
					Code:    500,
					Message: "Internal Server Error",
					Data:    err.Error(),
				})
			}
		}()

		return c.Next()
	})

	Group := app.Group("/api")
	Group.Route("/auth", AuthRoute)
	Group.Route("/usermanagement", UserRoute)
	Group.Route("/usertypemanagement", UserTypeRoute)
	Group.Route("/customer", CustomerRoute)
	Group.Route("/projectmanagement", ProjectRoute)
	Group.Route("/expensemanagement", ExpenseRoute)
	Group.Route("/dashboard", DashboardRoute)
}
