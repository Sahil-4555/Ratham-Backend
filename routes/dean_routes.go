package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github/Sahil-4555/ratham-backend/controllers"
)

func DeanRoutes(app *fiber.App) {
	app.Post("/dean/register", controllers.Register_Dean)
	app.Post("/dean/login", controllers.Login_Dean)
	app.Get("/dean/user", controllers.Dean)
	app.Post("/dean/logout", controllers.Logout_Dean)
}