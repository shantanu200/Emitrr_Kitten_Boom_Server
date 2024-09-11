package internals

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client


func InitRedis() {
	serviceURI := os.Getenv("REDIS_URL")

	if serviceURI == "" {
		panic("REDIS_URL is not set in environment")
	}

	addr, err := redis.ParseURL(serviceURI);

	if err != nil {
		panic(err)
	}


	log.Println("REDIS Connected to server");
	RDB = redis.NewClient(addr);

}