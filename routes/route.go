package routes

import (

	// import other routes...

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	AuthGroup := app.Group("/api")
	AuthGroup.Route("/auth", AuthRoute)
}
