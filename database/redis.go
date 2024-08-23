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


/*REDiS kullanırken dikkat edilecek yol
1-go get ile yükle github.com/go-redis/redis/v8

2- database klasöründe oluştur clean kod adına

3- Router dosyasında hem 3. parametre olarak hem de CRUD işlemlerine parametre olarak ekle

4- main.go sayfasında router fonksiyonunu çağırırken redis bağlantısını da parametre olarak ekle

5- handlers klaösründe ilgiliCRUD işlemlerine kodlarını yaz. 














*/
