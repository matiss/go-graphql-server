package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/graph-gophers/graphql-go"

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

func (h *GraphQLHandler) Query(c *fiber.Ctx) {
	ctx := h.ctx

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := c.BodyParser(&params); err != nil {
		fmt.Println(err)
		return
	}

	// Extract JWT token from Header and authenticate
	ctx = utils.SetAuth(ctx, c.Get("Authorization"), h.jwtSecret)

	// Set context values
	ctx = context.WithValue(ctx, "IP", c.IP())

	// Handle query
	response := h.schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	c.JSON(response)
}
