package tokens

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func DecodeToken(tokenString string, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil

}

func GetClientName(claims map[string]interface{}) (string, error) {
	sub, ok := claims["sub"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Invalid token body")
	}

	clientName, ok := sub["client_name"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid token body")
	}

	return clientName, nil
}
