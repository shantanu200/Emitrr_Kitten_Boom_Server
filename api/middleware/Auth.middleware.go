package middleware

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"path/filepath"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func LoadPublicKey() (*rsa.PublicKey, error) {
	keyPath, err := filepath.Abs("./../certs/public.pem")
	if err != nil {
		return nil, err
	}

	publicKeyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}



func GetJWTConfig() fiber.Handler {
	publicKey, err := LoadPublicKey()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Public Key Pulled successfully")

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    publicKey,
		},
	})

}
