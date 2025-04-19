package storage

import (
    "github.com/redis/go-redis/v9"
)

var RedisErr = redis.Nil

var RedisDB = redis.NewClient(&redis.Options{
    Addr: "localhost:6379", // Redis address
})

