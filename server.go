package main

import (
	"PBD_backend_go/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func runEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {

	app := fiber.New()

	// Define your routes here
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	runEnv()

	routes.SetupRoutes(app)

	app.Listen(`127.0.0.1:` + os.Getenv(`GO_PORT`))
}
