package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/storage"
)

type Rate_token struct {
	tokens        float64
	maxTokens     float64
	refilRate     float64
	lastRefilTime time.Time
}

// called on signup
func NewRateToken(ctx context.Context, jti string, expiration time.Duration, maxToken, refilRate float64) error {
	// first check if there is a rate token for an access token
	_, err := storage.RedisDB.HGet(ctx, "rate_token:"+jti, "tokens").Result()
	if err != nil && err != storage.RedisErr { // an error occured when getting the rate token
		return err
	} else if err == storage.RedisErr { // there is no rate token
		err := storage.RedisDB.HSet(ctx, "rate_token:"+jti,
			Rate_token{
				tokens:        maxToken - 1,
				maxTokens:     maxToken,
				refilRate:     refilRate,
				lastRefilTime: time.Now(),
			}, expiration).Err()
		if err != nil {
			return err
		}
		err = storage.RedisDB.Expire(ctx, jti, expiration).Err()
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("user already exists and has been given a rate token")
}

// called on every request
func CheckToken(ctx context.Context, jti string) (bool, error) {
	token, err := storage.RedisDB.HGet(ctx, "rate_token:"+jti, "tokens").Result()
	if err != nil {
		return false, err
	}

	err = refillToken(ctx, jti)
	if err != nil {
		return false, err
	}

	parsedToken, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return false, err
	}

	if parsedToken < 1 {
		return false, nil
	}

	parsedToken--
	err = storage.RedisDB.HSet(ctx, "rate_token:"+jti, "tokens", parsedToken).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func refillToken(ctx context.Context, jti string) error {
	rn := time.Now()

	token, err := storage.RedisDB.HGetAll(ctx, "rate_token:"+jti).Result()
	if err != nil {
		return err
	}
	// parsing fields to float from string
	lastRefilTime, err := time.Parse(time.RFC3339, token["lastRefilTime"])
	if err != nil {
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

	err = storage.RedisDB.HSet(ctx, "rate_token:"+jti, "tokens", tokenCount, "lastRefilTime", rn).Err()
	if err != nil {
		return err
	}

	return nil
}
