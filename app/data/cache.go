package data

import (
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

var cache *redis.Client

func Cache() *redis.Client {
	if cache == nil {
		d, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)

		if err != nil {
			log.Fatal(err)
		}

		cache = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       int(d),
		})
	}

	return cache
}
