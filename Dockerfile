FROM golang:alpine AS build
WORKDIR /build/
COPY . .
RUN go mod download
RUN go get -u github.com/go-bindata/go-bindata/...
RUN go generate ./schema
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -v' -v -o ./dist/api_server ./cmd/server/*.go

FROM scratch
WORKDIR /
COPY --from=build /build/dist/api_server /
COPY --from=build /build/config/server.toml /
ENTRYPOINT ["./api_server"]
EXPOSE 3035