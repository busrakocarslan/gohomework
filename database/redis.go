package database

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func UseRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis'e bağlanılamadı: %v", err)
	}

	return rdb

}