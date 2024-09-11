package internals

import (
	"fmt"
	"kitten-server/constants"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client


func InitRedis() {
	serviceURI := constants.GetENVConstant("REDIS_URL")

	if serviceURI == "" {
		fmt.Println("REDIS_URL is not set in environment")

		serviceURI = os.Getenv("REDIS_URL")
	}

	addr, err := redis.ParseURL(serviceURI);

	if err != nil {
		panic(err)
	}


	log.Println("REDIS Connected to server");
	RDB = redis.NewClient(addr);

}