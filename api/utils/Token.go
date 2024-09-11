package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
)



func GenerateAccessToken(userName string) (string, error) {
	keyPath,err := filepath.Abs("/app/certs/private.pem")
	
	if err != nil {
		return "",err
	}
	privateKey,err := os.ReadFile(keyPath);

	if err != nil {
		return "",err
	}

	rsapPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": userName,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(rsapPrivateKey)

	if err != nil {
		return "", nil
	}

	return tokenString, nil
}