.PHONY: build build-migrations run test migrate schema coverage fmt doc

ifndef VERBOSE
MAKEFLAGS+=--no-print-directory
endif

ifeq ($(UNAME),Darwin)
ECHO=echo
else
ECHO=echo -e
endif

# Package
PACKAGE_NAME=go-graphql-server
PACKAGE_VERSION=0.0.1-alpha
BUILD=$(shell git rev-list --count HEAD)
ARCHITECTURE=amd64
# LDFLAGS=-ldflags '-v'
LDFLAGS=-ldflags '-w -s -v'

SRCS=./cmd/server/*.go
SRCS_MIGRATIONS=./cmd/migrations/*.go

default: build

build:
	-@$(ECHO) "\n\033[0;35m%%% Building libraries and tools\033[0m"
	-@$(ECHO) "Building..."
	CGO_ENABLED=0 go build $(LDFLAGS) -v -o ./dist/$(PACKAGE_NAME) $(SRCS)
	-@$(ECHO) "\n\033[1;32mCONGRATULATIONS COMRADE!\033[0;32m\nDone!\033[0m\n"

build-migrations:
	-@$(ECHO) "\n\033[0;35m%%% Building migrations\033[0m"
	rm -rf ./dist/migrations
	-@$(ECHO) "Building migrations"
	CGO_ENABLED=0 go build $(LDFLAGS) -v -o ./dist/migrations $(SRCS_MIGRATIONS)
	-@$(ECHO) "\n\033[1;32mCONGRATULATIONS COMRADE!\033[0;32m\nDone!\033[0m\n"

run:
	go run ./cmd/server/main.go -c ./config/server.toml -D

test:
	-@$(ECHO) "\n\033[0;35m%%% Running tests\033[0m"
	go test -v ./...

coverage:
	-@$(ECHO) "\n\033[0;35m%%% Running test coverage\033[0m"
	go test -cover ./...

migrate:
	go run ./cmd/migrations/main.go migrate -c ./config/server.toml -D

schema:
	go generate ./schema

doc:
  godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./...