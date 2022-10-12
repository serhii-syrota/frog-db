package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/ssyrota/frog-db/src/core/db"
	"github.com/ssyrota/frog-db/src/core/db/dbtypes"
	"github.com/ssyrota/frog-db/src/core/db/schema"
	"github.com/ssyrota/frog-db/src/core/db/table"
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
	server.RegisterHandlers(
		r.Group("", middleware.OapiRequestValidator(swaggerFile)),
		server.NewStrictHandler(&handler{s.db}, []server.StrictMiddlewareFunc{}))

	swagger, err := server.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %w", err)
	}
	err = RegisterSwaggerHandler(r, swagger)
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
	tableSchema := schema.T{}
	for _, s := range *request.Body.Schema {
		tableSchema[s.Column] = dbtypes.Type(s.Type)
	}
	res, err := h.db.Execute(&db.CommandCreateTable{Name: *request.Body.TableName, Schema: tableSchema})
	if err != nil {
		return server.CreateTabledefaultJSONResponse{Body: server.Error{Message: err.Error()}, StatusCode: http.StatusConflict}, nil
	}
	message, ok := (*res)[0]["message"].(string)
	if !ok {
		return server.CreateTabledefaultJSONResponse{Body: server.Error{Message: "failed to create table"}, StatusCode: http.StatusInternalServerError}, nil
	}

	return server.CreateTable200JSONResponse{Message: message}, nil
}

// DeleteTable implementation.
func (h *handler) DeleteTable(ctx context.Context, request server.DeleteTableRequestObject) (server.DeleteTableResponseObject, error) {
	res, err := h.db.Execute(&db.CommandDropTable{Name: request.Name})
	if err != nil {
		return server.DeleteTabledefaultJSONResponse{Body: server.Error{Message: err.Error()}, StatusCode: http.StatusConflict}, nil
	}
	message, ok := (*res)[0]["message"].(string)
	if !ok {
		return server.DeleteTabledefaultJSONResponse{Body: server.Error{Message: "failed to drop table"}, StatusCode: http.StatusInternalServerError}, nil
	}

	return server.DeleteTable200JSONResponse{Message: message}, nil
}

// DeleteDuplicateRows implementation.
func (h *handler) DeleteDuplicateRows(ctx context.Context, request server.DeleteDuplicateRowsRequestObject) (server.DeleteDuplicateRowsResponseObject, error) {
	res, err := h.db.Execute(&db.CommandRemoveDuplicates{From: request.Name})
	if err != nil {
		return server.DeleteDuplicateRowsdefaultJSONResponse{Body: server.Error{Message: err.Error()}, StatusCode: http.StatusConflict}, nil
	}
	message, ok := (*res)[0]["message"].(string)
	if !ok {
		return server.DeleteDuplicateRowsdefaultJSONResponse{Body: server.Error{Message: "failed to drop table"}, StatusCode: http.StatusInternalServerError}, nil
	}

	return server.DeleteDuplicateRows200JSONResponse{Message: message}, nil
}

// DeleteRows implementation.
func (h *handler) DeleteRows(ctx context.Context, request server.DeleteRowsRequestObject) (server.DeleteRowsResponseObject, error) {
	res, err := h.db.Execute(&db.CommandDelete{From: request.Name, Conditions: RowToColumnSet(*request.Body)})
	if err != nil {
		return server.DeleteRowsdefaultJSONResponse{Body: server.Error{Message: err.Error()}, StatusCode: http.StatusConflict}, nil
	}
	message, ok := (*res)[0]["message"].(string)
	if !ok {
		return server.DeleteRowsdefaultJSONResponse{Body: server.Error{Message: "failed to delete rows"}, StatusCode: http.StatusInternalServerError}, nil
	}

	return server.DeleteRows200JSONResponse{Message: message}, nil
}

// InsertRows implementation.
func (h *handler) InsertRows(ctx context.Context, request server.InsertRowsRequestObject) (server.InsertRowsResponseObject, error) {
	data := make([]table.ColumnSet, len(*request.Body))
	for i, v := range *request.Body {
		data[i] = RowToColumnSet(v)
	}
	res, err := h.db.Execute(&db.CommandInsert{To: request.Name, Data: &data})
	if err != nil {
		return server.InsertRowsdefaultJSONResponse{Body: server.Error{Message: err.Error()}, StatusCode: http.StatusConflict}, nil
	}
	message, ok := (*res)[0]["message"].(string)
	if !ok {
		return server.InsertRowsdefaultJSONResponse{Body: server.Error{Message: "failed to insert rows"}, StatusCode: http.StatusInternalServerError}, nil
	}
	return server.InsertRows200JSONResponse{Message: message}, nil
}

// SelectRows implementation.
func (h *handler) SelectRows(ctx context.Context, request server.SelectRowsRequestObject) (server.SelectRowsResponseObject, error) {
	columns := request.Body.Columns
	conditions := RowToColumnSet(request.Body.Conditions)
	res, err := h.db.Execute(&db.CommandSelect{From: request.Name, Conditions: conditions, Fields: &columns})
	if err != nil {
		return server.SelectRowsdefaultJSONResponse{Body: server.Error{Message: err.Error()}, StatusCode: http.StatusConflict}, nil
	}
	response := make(server.SelectRows200JSONResponse, len(*res))
	for i, val := range *res {
		response[i] = ColumnSetToRows(val)
	}
	return response, nil
}

// UpdateRows implementation.
func (h *handler) UpdateRows(ctx context.Context, request server.UpdateRowsRequestObject) (server.UpdateRowsResponseObject, error) {
	conditions := RowToColumnSet(request.Body.Conditions)
	data := RowToColumnSet(request.Body.Data)
	res, err := h.db.Execute(&db.CommandUpdate{TableName: request.Name, Conditions: conditions, Data: data})
	if err != nil {
		return server.UpdateRowsdefaultJSONResponse{Body: server.Error{Message: err.Error()}, StatusCode: http.StatusConflict}, nil
	}
	message, ok := (*res)[0]["message"].(string)
	if !ok {
		return server.UpdateRowsdefaultJSONResponse{Body: server.Error{Message: "failed to update rows"}, StatusCode: http.StatusInternalServerError}, nil
	}
	return server.UpdateRows200JSONResponse{Message: message}, nil
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
		// Prevent range value pointer reference bug
		tableNameCopy := tableName
		res = append(res, server.TableSchema{TableName: &tableNameCopy, Schema: &schema})
	}
	return res, nil
}

func RowToColumnSet(row server.Row) table.ColumnSet {
	res := table.ColumnSet{}
	for k, v := range row {
		res[k] = v
	}
	return res
}
func ColumnSetToRows(row table.ColumnSet) server.Row {
	res := server.Row{}
	for k, v := range row {
		res[k] = v
	}
	return res
}
