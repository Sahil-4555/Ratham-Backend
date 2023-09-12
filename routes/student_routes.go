package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github/Sahil-4555/ratham-backend/controllers"
)

func StudentRoutes(app *fiber.App) {
	app.Post("/student/register", controllers.Register_Student)
	app.Post("/student/login", controllers.Login_Student)
	app.Get("/student/user", controllers.Student)
	app.Post("/student/logout", controllers.Logout_Student)
}