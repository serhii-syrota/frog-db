package web

import (
	"context"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/ssyrota/frog-db/src/web/server"
)

func Run() error {
	r := echo.New()
	swaggerFile, err := server.GetSwagger()
	if err != nil {
		return err
	}
	r.Use(middleware.OapiRequestValidator(swaggerFile))
	return nil
}

var _ server.StrictServerInterface = new(Handler)

type Handler struct {
}

// DbSchema implementation.
func (h *Handler) DbSchema(ctx context.Context, request server.DbSchemaRequestObject) (server.DbSchemaResponseObject, error) {
	return nil, nil
}
