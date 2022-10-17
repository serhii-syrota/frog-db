package web

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/tj/go/env"
)

func RegisterSwaggerHandler(r *echo.Echo, swagger *openapi3.T) error {
	swaggerMarshaled, err := swagger.MarshalJSON()
	if err != nil {
		return fmt.Errorf("swagger marshal: %w", err)
	}
	swaggerPage := SwaggerStaticHtml()

	r.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, swaggerPage)
	})
	r.GET("/docs/.json", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, swaggerMarshaled)
	})
	return nil
}
func SwaggerStaticHtml() string {
	host := "https://frogdb.herokuapp.com"
	if env.Get("ENV") == "DEV" {
		host = "http://localhost:8080"
	}

	return `<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="utf-8" />
	  <meta name="viewport" content="width=device-width, initial-scale=1" />
	  <meta
		name="description"
		content="SwaggerIU"
	  />
	  <title>SwaggerUI</title>
	  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.4.1/swagger-ui.css" />
	</head>
	<body>
	  <div id="swagger-ui"></div>
	  <script src="https://unpkg.com/swagger-ui-dist@4.4.1/swagger-ui-bundle.js" crossorigin></script>
	  <script src="https://unpkg.com/swagger-ui-dist@4.4.1/swagger-ui-standalone-preset.js" crossorigin></script>
	  <script>
		window.onload = () => {
		  window.ui = SwaggerUIBundle({
			url: "` + host + `/docs/.json",
			dom_id: '#swagger-ui',
		  });
		};
	  </script>
	</body>
	</html>`
}
