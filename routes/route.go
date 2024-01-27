package routes

import (

	// import other routes...

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	f := fiber.New()
	f.Use(logger.New())
	f.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	Group := app.Group("/api")
	Group.Route("/auth", AuthRoute)
	Group.Route("/usermanagement", UserManagementRoute)
}
