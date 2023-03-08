package tokens

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func DecodeToken(tokenString string, JWTPublicKey string) (map[string]interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(JWTPublicKey))
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid Token")
	}

	return claims, nil

}

func GetClientID(claims map[string]interface{}) (uuid.UUID, error) {
	sub, ok := claims["sub"].(map[string]interface{})
	if !ok {
		return uuid.UUID{}, fmt.Errorf("Invalid token body")
	}

	ID, ok := sub["client_id"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("Invalid token body")
	}

	clientID, err := uuid.Parse(ID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return clientID, nil
}
