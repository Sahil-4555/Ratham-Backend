package main

import (
	"github/Sahil-4555/ratham-backend/routes"
	"github/Sahil-4555/ratham-backend/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New();

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	configs.ConnectDB()
	routes.StudentRoutes(app);
	routes.DeanRoutes(app);
	routes.SessionRoutes(app);
	app.Listen(":8080")	
}