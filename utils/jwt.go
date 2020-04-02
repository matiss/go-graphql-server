package utils

import (
	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(exp int64, secret []byte, userId int32, authLevel AuthLevel) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	claims := t.Claims.(jwt.MapClaims)
	claims["sub"] = userId
	claims["exp"] = exp
	claims["auth"] = int(authLevel)

	tokenString, err := t.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(key *[]byte, msg string) (*jwt.Token, error) {
	return jwt.Parse(msg, func(token *jwt.Token) (interface{}, error) {
		return *key, nil
	})
}
