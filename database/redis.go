package database

import (
	"context"
	"log"
	"os"// çevresel değişkenler için 

	"github.com/go-redis/redis/v8"
)

func UseRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("Redis URL is not set in environment variables")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}
	rdb := redis.NewClient(opt)
	
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return rdb

}


/*REDiS kullanırken dikkat edilecek yol
1-go get ile yükle github.com/go-redis/redis/v8

2- database klasöründe oluştur clean kod adına

3- Router dosyasında hem 3. parametre olarak hem de CRUD işlemlerine parametre olarak ekle

4- main.go sayfasında router fonksiyonunu çağırırken redis bağlantısını da parametre olarak ekle

5- handlers klaösründe ilgiliCRUD işlemlerine kodlarını yaz. 














*/
