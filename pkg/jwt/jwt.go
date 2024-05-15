package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func getSigningKey() ([]byte, error) {
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))
	if len(mySigningKey) == 0 {
		return nil, fmt.Errorf("JWT secret key not set")
	}
	return mySigningKey, nil
}

func GenerateToken(userID string) (string, error) {
	mySigningKey, err := getSigningKey()
	if err != nil {
		return "", err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	mySigningKey, err := getSigningKey()
	if err != nil {
		return "", err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return mySigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	} else {
		return "", fmt.Errorf("invalid or expired token")
	}
}
