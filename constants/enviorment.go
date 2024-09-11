package constants

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetENVConstant(key string) string {
	if err := godotenv.Load("./../.env") ; err != nil {
		fmt.Println("Error loading .env file")
	}

	return os.Getenv(key)
	
}