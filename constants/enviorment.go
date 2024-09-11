package constants

import (
	"os"

	"github.com/joho/godotenv"
)

func GetENVConstant(key string) string {
	
	if err := godotenv.Load("./../.env") ; err != nil {
		panic(err)
	}
	
	return os.Getenv(key)
	
}