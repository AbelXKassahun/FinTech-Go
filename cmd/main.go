package main

import (
	"context"
	"log"
	"net/http"
	"os"


	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/api"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/storage"
)

// "github.com/AbelXKassahun/Digital-Wallet-Platform/internal/auth"


func main() {
	port := os.Getenv("APP_PORT")
	ctx := context.Background()

	// connecting to postgreSQL
	if err := storage.InitPostgres(); err != nil {
		storage.DB.Close()
		log.Fatal("Error connecting to Postgres:", err)
	}
	log.Println("Connected to Postgres")
	defer storage.DB.Close()

	// connecting to redis
	storage.InitRedis()
	if _, err := storage.RedisDB.Ping(ctx).Result(); err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	log.Println("Connected to Redis")

	log.Printf("Listening at port 8080 \n")
	if err := http.ListenAndServe(":"+port, api.Routes()); err != nil {
		panic(err)
	}

	// err := storage.RedisDB.HSet(context.Background(), "rate_token:1", map[string]interface{}{
	// 	"tokens":    2,
	// 	"MaxTokens": 3,
	// 	"RefilRate": 0.5,
	// 	// "LastRefilTime": time.Now(),
	// }).Err()

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// token, err := storage.RedisDB.HGet(context.Background(), "rate_token:1", "tokens").Result()
	// fmt.Println("token: ", token)
	// if err != nil {
	// 	fmt.Println("here")
	// }
}
