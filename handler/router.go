package handler

import (
	"pt_search_hos/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, staff *StaffHandler, patient *PatientHandler, jwtSecret string, isBlacklisted func(string) bool) {
	auth := middleware.JWTProtected(jwtSecret, isBlacklisted)

	// System
	app.Get("/hello", HelloHandler)

	// Staff (auth)
	app.Post("/staff/login", staff.Login)
	app.Post("/staff/create",  staff.Create)
	app.Get("/staff/hello", auth, staff.Hello)
	app.Get("/staff/logout", auth, staff.Logout)

	// Patient (auth)
	app.Post("/patient/search", auth, patient.Search)
	app.Get("/patient/search/:id", auth, patient.GetByID)
}
