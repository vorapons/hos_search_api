package ginhandler

import (
	ginmw "pt_search_hos/middleware/gin"

	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, staff *StaffHandler, patient *PatientHandler, jwtSecret string, isBlacklisted func(string) bool) {
	auth      := ginmw.JWTProtected(jwtSecret, isBlacklisted)
	authLimit := ginmw.AuthRateLimit() // 10 req/min — brute-force protection
	apiLimit  := ginmw.APIRateLimit()  // 100 req/min — general API protection

	// System
	r.GET("/",             func(c *gin.Context) { c.String(http.StatusOK, "Hospital Search API is running") })
	r.GET("/hello",        HelloHandler)
	r.GET("/swagger",      SwaggerHandler)
	r.GET("/openapi.yaml", OpenAPIHandler)

	// Staff — public (strict rate limit)
	r.POST("/staff/login",  authLimit, staff.Login)
	r.POST("/staff/create", authLimit, staff.Create)

	// Staff — authenticated
	r.GET("/staff/hello",  apiLimit, auth, staff.Hello)
	r.POST("/staff/logout", apiLimit, auth, staff.Logout)

	// Patient — authenticated
	r.POST("/patient/search",    apiLimit, auth, patient.Search)
	r.GET("/patient/search/:id", apiLimit, auth, patient.GetByID)
}
