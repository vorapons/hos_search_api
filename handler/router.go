package handler

import (
	"pt_search_hos/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, staff *StaffHandler, patient *PatientHandler, jwtSecret string, isBlacklisted func(string) bool) {
	auth      := middleware.JWTProtected(jwtSecret, isBlacklisted)
	authLimit := middleware.AuthRateLimit() // 10 req/min — brute-force protection
	apiLimit  := middleware.APIRateLimit()  // 100 req/min — general API protection

	// System
	app.Get("/hello", HelloHandler)

	// Staff — public (strict rate limit)
	app.Post("/staff/login",  authLimit, staff.Login)
	app.Post("/staff/create", authLimit, staff.Create)

	// Staff — authenticated
	app.Get("/staff/hello",  apiLimit, auth, staff.Hello)
	app.Get("/staff/logout", apiLimit, auth, staff.Logout)

	// Patient — authenticated
	app.Post("/patient/search",    apiLimit, auth, patient.Search)
	app.Get("/patient/search/:id", apiLimit, auth, patient.GetByID)
}
