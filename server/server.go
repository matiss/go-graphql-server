package server

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/limiter"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"

	"github.com/matiss/go-graphql-server/resolvers"
	"github.com/matiss/go-graphql-server/schema"
	"github.com/matiss/go-graphql-server/services"
)

var (
	JWTSecret = []byte("4S72scrC7ESfJoyMST4EhnF2CyrvA0Xc79JaP2MU2onLTxGmsrwuZRP6X5zDIJA")
)

func Run(configPath string) {
	// Load config
	config := services.ConfigService{}
	err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Connect to Postgres database
	pgService := services.PGService{}
	err = pgService.Connect(config.PG.Address, config.PG.User, config.PG.Password, config.PG.Database, config.PG.PoolSize)
	if err != nil {
		panic(err)
	}

	userService := services.NewUserService(pgService.DB)

	s := fiber.New()

	// Set prefork
	s.Settings.Prefork = false

	// Recover middleware
	s.Use(recover.New(recover.Config{
		// Config is optional
		Handler: func(c *fiber.Ctx, err error) {
			c.SendString(err.Error())
			c.SendStatus(500)
		},
	}))

	// Create a rate limiter struct.
	rateLimiter := limiter.Config{
		Timeout: 1,
		Max:     config.HTTP.RateLimit,
	}
	s.Use(limiter.New(rateLimiter))

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"content-type", "Authorization"},
	}
	s.Use(cors.New(corsConfig))

	// Logger middleware
	s.Use(logger.New())

	// Root handler
	s.Get("/", RootHandler)

	// GraphQL resolver
	reolver := &resolvers.Resolver{
		JWTSecret:     JWTSecret,
		TokenTTL:      config.Auth.TokenTTL,
		TokenTTLRenew: config.Auth.TokenTTLRenew,
		UserService:   userService,
	}

	// GraphQL schema
	schema := graphql.MustParseSchema(schema.GetRootSchema(), reolver)

	// GraphQL Handler
	graphqlHandler := NewGraphQLHandler(ctx, &JWTSecret, schema)

	// GraphQL handler
	s.Post("/query", graphqlHandler.Query)

	// Handle robots.txt file
	s.Get("/robots.txt", RobotsTXTHandler)

	// Start server
	s.Listen(config.HTTP.Address)
}
