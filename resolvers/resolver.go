package resolvers

import (
	"github.com/matiss/go-graphql-server/services"
)

// Resolver is the root resolver
type Resolver struct {
	JWTSecret     []byte
	TokenTTL      int64
	TokenTTLRenew int64

	UserService *services.UserService
}
