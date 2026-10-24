package router

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var swaggerFiles embed.FS

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Visa Assignment API - Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: "/swagger/openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true
      });
    };
  </script>
</body>
</html>`

func registerSwaggerRoutes(r *gin.Engine) {
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})
	r.GET("/swagger/openapi.yaml", func(c *gin.Context) {
		spec, err := swaggerFiles.ReadFile("openapi.yaml")
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to load OpenAPI specification")
			return
		}
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", spec)
	})
}
