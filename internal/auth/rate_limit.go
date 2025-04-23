package auth

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/storage"
)

type Rate_token struct {
	Tokens        float64
	MaxTokens     float64
	RefilRate     float64
	LastRefilTime time.Time
}

// called on signup
func NewRateToken(ctx context.Context, key string, expiration time.Duration, maxToken, refilRate float64) error {
	// first check if there is a rate token for an access token
	val, err := storage.RedisDB.HGet(ctx, "rate_token:"+key, "tokens").Result()
	if err != nil && err != storage.RedisErr { // an error occured when getting the rate token
		return err
	} else if err == storage.RedisErr && val == "" { // there is no rate token

		rateMap := map[string]interface{}{
			"tokens":          maxToken - 1,
			"maxTokens":      maxToken,
			"refilRate":      refilRate,
			"lastRefilTime": time.Now().Format(time.RFC3339),
		}
		
		_, err := storage.RedisDB.HSet(ctx, "rate_token:"+key, rateMap).Result()
		if err != nil {
			return err
		}

		wasSet, err := storage.RedisDB.Expire(ctx, "rate_token:"+key, expiration).Result()
		if err != nil {
			log.Println("error expiring redis hash")
			return err
		}

		if wasSet {
			log.Printf("Expiration set for key '%s' to %s\n", "rate_token:"+key, expiration)
		} else {
			log.Printf("Expiration could not be set for key '%s' (key might not exist?)\n", "rate_token:"+key)
		}
	
		// checking the TTL (Time To Live)
		ttl, err := storage.RedisDB.TTL(ctx, "rate_token:"+key).Result()
		if err != nil {
			log.Println("Error getting TTL:", err)
		} else {
			if ttl < 0 {
				// TTL < 0 indicates no expiry (-1) or key doesn't exist (-2)
				log.Printf("Key '%s' has no expiration (TTL: %s)\n", "rate_token:"+key, ttl)
			} else {
				log.Printf("Key '%s' will expire in approximately %s\n", "rate_token:"+key, ttl)
			}
		}
		
		return nil
	}

	return fmt.Errorf("user already exists and has been given a rate token")
}

// called on every request
func CheckToken(ctx context.Context, key string) (bool, error) {
	token, err := storage.RedisDB.HGet(ctx, "rate_token:"+key, "tokens").Result()
	fmt.Println("token: ", token)

	if err != nil && err != storage.RedisErr {
		return false, err
	} else if err == storage.RedisErr {
		return false, storage.RedisErr
	}

	err = refillToken(ctx, key)
	if err != nil {
		return false, err
	}

	parsedToken, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return false, err
	}

	if parsedToken < 1 { // true false (the user has exceeded the rate limit)
		return false, nil
	}

	parsedToken--
	err = storage.RedisDB.HSet(ctx, "rate_token:"+key, "tokens", parsedToken).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func refillToken(ctx context.Context, key string) error {
	rn := time.Now()

	token, err := storage.RedisDB.HGetAll(ctx, "rate_token:"+key).Result()
	if err != nil {
		return err
	}
	fmt.Println("rateMap", token)
	// parsing fields to float from string
	lastRefilTime, err := time.Parse(time.RFC3339, token["lastRefilTime"])
	if err != nil {
		fmt.Println("here time")
		return err
	}
	refileRate, err := strconv.ParseFloat(token["refilRate"], 64)
	if err != nil {
		return err
	}
	tokenCount, err := strconv.ParseFloat(token["tokens"], 64)
	if err != nil {
		return err
	}
	maxTokenCount, err := strconv.ParseFloat(token["maxTokens"], 64)
	if err != nil {
		return err
	}

	elapsed := rn.Sub(lastRefilTime).Seconds()
	tokenCount += refileRate * elapsed

	if tokenCount > maxTokenCount {
		tokenCount = maxTokenCount
	}

	err = storage.RedisDB.HSet(ctx, "rate_token:"+key, "tokens", tokenCount, "lastRefilTime", rn).Err()
	if err != nil {
		return err
	}

	return nil
}
