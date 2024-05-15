package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateToken(userID string) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))
	if len(mySigningKey) == 0 {
		return "", fmt.Errorf("JWT secret key not set")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))
	if len(mySigningKey) == 0 {
		return "", fmt.Errorf("JWT secret key not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}
