package data

import (
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

var redisClient *redis.Client

func Redis() *redis.Client {
	if redisClient == nil {
		d, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)

		if err != nil {
			log.Fatal(err)
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       int(d),
		})
	}

	return redisClient
}
