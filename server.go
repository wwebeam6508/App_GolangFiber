package main

import (
	"context"

	"PBD_backend_go/configuration"

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

	// call the ConnectToMongoDB function
	client, err := configuration.ConnectToMongoDB()
	if err != nil {
		panic(err)
	}

	// Close the connection
	defer client.Disconnect(context.Background())

	app.Listen(":3000")
}
