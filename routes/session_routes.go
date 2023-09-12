package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github/Sahil-4555/ratham-backend/controllers"
)

func SessionRoutes(app *fiber.App) {
	app.Post("/session/addsession", controllers.AddNewSession)
	app.Get("/session/getfreesession", controllers.GetAllFreeSessions)
	app.Post("/session/booksession/:id", controllers.BookASession)
	app.Get("/session/getupcomingfreesession", controllers.GetUpcomingFreeSessions)
}	