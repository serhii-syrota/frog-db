package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/ssyrota/frog-db/src/core/db"
	"github.com/ssyrota/frog-db/src/web/server"
)

func New(db *db.Database, port uint16) *WebServer {
	return &WebServer{port, db}
}

type WebServer struct {
	port uint16
	db   *db.Database
}

func (s *WebServer) Run() error {
	r := echo.New()
	swaggerFile, err := server.GetSwagger()
	if err != nil {
		return err
	}
	r.Use(echo_middleware.Logger(), echo_middleware.Recover())
	apiGroup := r.Group("", middleware.OapiRequestValidator(swaggerFile))
	r.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "string")
	})
	server.RegisterHandlers(apiGroup, server.NewStrictHandler(&handler{s.db}, []server.StrictMiddlewareFunc{}))

	swagger, err := server.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %w", err)
	}
	err = RegisterSwaggerHandler(r, "http://localhost:8080", swagger)
	if err != nil {
		return fmt.Errorf("swagger register: %w", err)
	}
	return r.Start(fmt.Sprintf(":%d", s.port))
}

type handler struct {
	db *db.Database
}

// Verify handler implements server.StrictServerInterface
var _ server.StrictServerInterface = new(handler)

// CreateTable implementation.
func (h *handler) CreateTable(ctx context.Context, request server.CreateTableRequestObject) (server.CreateTableResponseObject, error) {
	return nil, nil
}

// DeleteTable implementation.
func (h *handler) DeleteTable(ctx context.Context, request server.DeleteTableRequestObject) (server.DeleteTableResponseObject, error) {
	return nil, nil
}

// DeleteDuplicateRows implementation.
func (h *handler) DeleteDuplicateRows(ctx context.Context, request server.DeleteDuplicateRowsRequestObject) (server.DeleteDuplicateRowsResponseObject, error) {
	return nil, nil
}

// DeleteRows implementation.
func (h *handler) DeleteRows(ctx context.Context, request server.DeleteRowsRequestObject) (server.DeleteRowsResponseObject, error) {
	return nil, nil
}

// InsertRows implementation.
func (h *handler) InsertRows(ctx context.Context, request server.InsertRowsRequestObject) (server.InsertRowsResponseObject, error) {
	return nil, nil
}

// SelectRows implementation.
func (h *handler) SelectRows(ctx context.Context, request server.SelectRowsRequestObject) (server.SelectRowsResponseObject, error) {
	return nil, nil
}

// UpdateRows implementation.
func (h *handler) UpdateRows(ctx context.Context, request server.UpdateRowsRequestObject) (server.UpdateRowsResponseObject, error) {
	return nil, nil
}

// DbSchema implementation.
func (h *handler) DbSchema(ctx context.Context, request server.DbSchemaRequestObject) (server.DbSchemaResponseObject, error) {
	schema, err := h.db.IntrospectSchema()
	if err != nil {
		return nil, err
	}
	res := server.DbSchema200JSONResponse{}
	for tableName, tableSchema := range schema {
		schema := []server.Schema{}
		for column, columnType := range tableSchema {
			schema = append(schema, server.Schema{Column: column, Type: server.SchemaType(columnType)})
		}
		res = append(res, server.TableSchema{TableName: &tableName, Schema: &schema})
	}

	return res, nil
}
