package cache

import (
	"errors"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func RedisConnection() (*redis.Client, error) {
	dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	redisConnectionURL := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	options := &redis.Options{
		Addr:     redisConnectionURL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	}
	client := redis.NewClient(options)
	if client == nil {
		return nil, errors.New("redis connection not possible")
	}
	return client, nil
}
