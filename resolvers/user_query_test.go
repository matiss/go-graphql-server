package resolvers

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/gqltesting"

	"github.com/matiss/go-graphql-server/schema"
	"github.com/matiss/go-graphql-server/services"
	"github.com/matiss/go-graphql-server/utils"
)

var (
	rootSchema *graphql.Schema
	ctx        context.Context
)

func init() {
	// Load config
	config := services.ConfigService{}
	err := config.Load("../config/test.toml")
	if err != nil {
		panic(err)
	}

	// Connect to Postgres database
	pgService := services.PGService{}
	err = pgService.Connect(config.PG.Address, config.PG.User, config.PG.Password, config.PG.Database, config.PG.PoolSize)
	if err != nil {
		panic(err)
	}

	userService := services.NewUserService(pgService.DB)

	// GraphQL resolver
	reolver := &Resolver{
		JWTSecret:     []byte("testtest"),
		TokenTTL:      config.Auth.TokenTTL,
		TokenTTLRenew: config.Auth.TokenTTLRenew,
		UserService:   userService,
	}

	rootSchema = graphql.MustParseSchema(schema.GetRootSchema(), reolver)

	ctx = context.Background()
	ctx = context.WithValue(ctx, "IP", "127.0.0.1")
	ctx = context.WithValue(ctx, "auth", utils.AuthUser)
	ctx = context.WithValue(ctx, "userID", "adam@test.com")
}

func TestUser(t *testing.T) {
	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Context: ctx,
			Schema:  rootSchema,
			Query: `
        {
          user(email:"adam@test.com") {
            email
            name
            password
          }
        }
      `,
			ExpectedResult: `
        {
          "user": {
            "email": "adam@test.com",
            "name": "Adam Tester",
            "password": "********"
          }
        }
      `,
		},
	})
}
