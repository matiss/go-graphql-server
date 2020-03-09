package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matiss/go-graphql-server/server"
)

const version = "0.0.1"

var (
	// A path to the config file.
	flConfigFile string
	// Print version
	flVersion bool
	// Enable debug mode
	flDebug bool
)

func init() {
	flag.BoolVar(&flVersion, "v", false, "Print version information and quit")
	flag.BoolVar(&flDebug, "D", false, "Enable debug mode")
	flag.StringVar(&flConfigFile, "c", "./server.toml", "path to config file")
	flag.Parse()

	flConfigFile, _ = filepath.Abs(flConfigFile)
}

func main() {
	if flVersion {
		fmt.Printf(" ⇛ GraphQL API Server v%s", version)
		return
	}

	if flDebug {
		os.Setenv("DEBUG", "1")
	}

	server.Run(flConfigFile)
}
