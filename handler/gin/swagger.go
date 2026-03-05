package ginhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const swaggerHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>Hospital Patient Search API</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      SwaggerUIBundle({
        url:     "/openapi.yaml",
        dom_id:  "#swagger-ui",
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
        layout:  "BaseLayout",
        deepLinking: true,
      })
    </script>
  </body>
</html>`

// SwaggerHandler serves the Swagger UI at GET /swagger
func SwaggerHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, swaggerHTML)
}

// OpenAPIHandler serves the raw openapi.yaml at GET /openapi.yaml
func OpenAPIHandler(c *gin.Context) {
	c.File("./openapi.yaml")
}
