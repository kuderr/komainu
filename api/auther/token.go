package auther

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type TokenBody struct {
	ClientName string `json:"client_name"`
}

type DecodedToken struct {
	Subject TokenBody `json:"sub"`
	jwt.StandardClaims
}

func decodeToken(token string, secret string) (DecodedToken, error) {
	creds := strings.Replace(token, "Bearer ", "", 1)

	tk := &DecodedToken{}
	decoded, err := jwt.ParseWithClaims(creds, tk, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if !decoded.Valid {
		err = errors.New("Invalid Token")
	}

	return *tk, err
}
