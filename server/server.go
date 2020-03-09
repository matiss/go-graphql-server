package server

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_echo"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	// Create a rate limiter struct.
	rateLimiter := tollbooth.NewLimiter(config.HTTP.RateLimit, nil)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, echo.HeaderAccessControlMaxAge},
	}))

	// Root handler
	e.GET("/", RootHandler, tollbooth_echo.LimitHandler(rateLimiter))

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
	e.POST("/query", graphqlHandler.Query, tollbooth_echo.LimitHandler(rateLimiter))

	// Handle robots.txt file
	e.GET("/robots.txt", RobotsTXTHandler, tollbooth_echo.LimitHandler(rateLimiter))

	// Stop here if its Preflighted OPTIONS request
	e.OPTIONS("/*url", HTTPOptions)

	// Start server
	go func() {
		if err := e.Start(config.HTTP.Address); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	sigs := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})

	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Cleanup
		pgService.Kill()

		// Stop Server
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}

		close(cleanupDone)
	}()

	<-cleanupDone
}
