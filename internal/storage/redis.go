package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisDB *redis.Client
var RedisErr = redis.Nil

func InitRedis() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // if no password
		DB:       0,
	})
}

// this function returns an interface so you often need to cast it (valStr,ok := val.(string))
func DoCommand(ctx context.Context, args ...interface{}) (interface{}, error) {
    return RedisDB.Do(ctx, args...).Result()
}
/*
// Set a key
_, err := DoCommand(ctx, "SET", "mykey", "myvalue")
if err != nil {
    log.Fatal(err)
}

// Get a key
val, err := DoCommand(ctx, "GET", "mykey")
if err != nil {
    log.Fatal(err)
}
valStr, ok := val.(string)
if !ok {
    log.Println("unexpected type")
}
*/