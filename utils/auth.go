package utils

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type AuthLevel int

const (
	AuthAnonymous AuthLevel = iota
	AuthUser
	AuthAdmin
)

func (a AuthLevel) Authorized(level AuthLevel) bool {
	return (level >= a)
}

// IsAuthorized is a auth helper for resolvers
func IsAuthorized(ctx context.Context, requiredLevel AuthLevel) (bool, error) {
	// Bypass auth if resource is for public use
	if requiredLevel == AuthAnonymous {
		return true, nil
	}

	// Check access level
	if authLevel, ok := ctx.Value("auth").(AuthLevel); !ok || !requiredLevel.Authorized(authLevel) {
		return false, fmt.Errorf("Unauthorized")
	}

	return true, nil
}

func SetAuth(ctx context.Context, authToken string, jwtSecret *[]byte) context.Context {
	// Default access to public/anonymous
	authLevel := AuthAnonymous
	var userID int32 = 0

	if len(authToken) > 7 {
		authToken = authToken[7:]

		token, err := ParseJWT(jwtSecret, authToken)
		if err == nil && token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			// Get userID
			if ID, ok := claims["sub"].(float64); ok {
				userID = int32(ID)
			}

			// Set auth level
			if lvl, ok := claims["auth"].(float64); ok {
				authLevel = AuthLevel(lvl)
			}
		}
	}

	// Set context values
	ctx = context.WithValue(ctx, "auth", authLevel)
	ctx = context.WithValue(ctx, "userID", userID)

	return ctx
}
