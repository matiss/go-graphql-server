## Go GraphQL API Server

This project aims to use [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go) to build a GraphQL API server.

#### RoadMap:

- [x] Integrated with pg
- [x] Database migrations
- [x] Integrated with graphql-go
- [x] Use go-bindata to generate Go code from .graphql file
- [x] Add authentication & authorization
- [x] Add simple unit test cases
    
#### Requirement:

1. Postgres database
2. Golang
3. GNU Make (Optional)

#### Usage:

1. Create database in Postgres and update server.toml configuration and run migrations
    ```
    make migrate
    ```

2. Install go-bindata
    ```
    go get -u github.com/go-bindata/go-bindata...
    ```

3. Run the following command at root directory to generate Go code from .graphql file
    ```
    make schema
    ```
    There would be bindata.go generated under `schema` folder


4. Start the server (Ensure your postgres database is live and its setting in server.toml is correct)
    ```
    make run
    ```

#### Test:

- Run Unit Tests
    ```
    make test
    ```
