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

	// Create a new fiber instance with custom config
	app := fiber.New()

	runEnv()

	routes.SetupRoutes(app)

	app.Listen(`localhost:` + os.Getenv(`GO_PORT`))
}
