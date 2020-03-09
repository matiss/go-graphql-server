package server

import (
	"context"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo/v4"

	"github.com/matiss/go-graphql-server/utils"
)

type GraphQLHandler struct {
	ctx       context.Context
	jwtSecret *[]byte
	schema    *graphql.Schema
}

func NewGraphQLHandler(ctx context.Context, jwtSecret *[]byte, schema *graphql.Schema) *GraphQLHandler {
	handler := GraphQLHandler{
		ctx,
		jwtSecret,
		schema,
	}

	return &handler
}

func (h *GraphQLHandler) Query(c echo.Context) error {
	ctx := h.ctx

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := c.Bind(&params); err != nil {
		return err
	}

	// Extract JWT token from Header and authenticate
	ctx = utils.SetAuth(ctx, c.Request().Header.Get("Authorization"), h.jwtSecret)

	// Set context values
	ctx = context.WithValue(ctx, "IP", c.RealIP())

	// Handle query
	response := h.schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	return c.JSON(http.StatusOK, response)
}
