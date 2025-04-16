package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/storage"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// replace this in prod
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(payload map[string]interface{}, expiration time.Duration) (string, error) {
	jti := uuid.New().String()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload["exp"] = time.Now().Add(expiration).Unix()
	payload["jti"] = jti

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal header: %v", err)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal payload: %v", err)
	}

	headerEncoded := utils.Base64URLEncode(headerJSON)
	payloadEncoded := utils.Base64URLEncode(payloadJSON)

	unsigned := headerEncoded + "." + payloadEncoded
	signature := utils.Sign(unsigned, jwtSecret)

	return unsigned + "." + signature, nil
}

func VerifyJWT(w http.ResponseWriter, r *http.Request) bool {
	tokenStr := utils.GetJWTFromRequest(w, r)

	// verified_token, err := auth.VerifyJWT(token, r.Context())
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	
    if err != nil  {
        w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return false
    } else if !token.Valid {
		w.WriteHeader(401)
		w.Write([]byte("Token is invalid"))
		return false
	}

    // [To Do] - uncomment this code if you want to check if the token is blacklisted

    // claims := token.Claims.(jwt.MapClaims)
    // jti := claims["jti"].(string)

    // blacklisted, err := IsBlacklisted(jti, r.Context())
    // if blacklisted {
	// 	w.WriteHeader(401)
    //     w.Write([]byte("Token is invalid"))
	// 	return false
    // }else if err != nil && !blacklisted {
    //     w.WriteHeader(500)
    //     w.Write([]byte(err.Error()))
	// 	return false
    // } 
	
	w.Write([]byte("Token Verified"))
	return true
}



func Logout(tokenStr string, ctx context.Context) error {
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	claims := token.Claims.(jwt.MapClaims)
	jti := claims["jti"].(string)
	exp := time.Until(time.Unix(int64(claims["exp"].(float64)), 0))

	return BlacklistToken(ctx, jti, exp)
}

func BlacklistToken(ctx context.Context, jti string, exp time.Duration) error {
	return storage.RedisDB.Set(ctx, "blacklist:"+jti, "1", exp).Err()
}

func IsBlacklisted(jti string, ctx context.Context) (bool, error) {
	val, err := storage.RedisDB.Get(ctx, "blacklist:"+jti).Result()
	if err == storage.RedisErr { // not blacklisted
		return false, nil
	} else if err != nil { // redis error
		return false, err
	}
	return val == "1", nil
}