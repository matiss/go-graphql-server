package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	migrations "github.com/robinjoseph08/go-pg-migrations/v2"

	_ "github.com/matiss/go-graphql-server/migrations"
	"github.com/matiss/go-graphql-server/services"
)

const directory = "migrations"

var (
	// A path to the config file.
	flConfigFile string
	// Enable debug mode
	flDebug bool
)

func init() {
	flag.BoolVar(&flDebug, "D", false, "Enable debug mode")
	flag.StringVar(&flConfigFile, "c", "./config/server.toml", "path to config file")
	flag.Parse()

	flConfigFile, _ = filepath.Abs(flConfigFile)
}

func main() {
	if flDebug {
		os.Setenv("DEBUG", "1")
	}

	// Load config
	config := services.ConfigService{}
	err := config.Load(flConfigFile)
	if err != nil {
		panic(err)
	}

	// Connect to Postgres database
	pgService := services.PGService{}
	err = pgService.Connect(config.PG.Address, config.PG.User, config.PG.Password, config.PG.Database, config.PG.PoolSize)
	if err != nil {
		panic(err)
	}

	err = migrations.Run(pgService.DB, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
