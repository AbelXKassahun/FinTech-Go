package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/storage"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//  TODO: store the jti in redis

type Claims struct {
	UserID string `json:"user_id"`
	Tier  string `json:"tier"`
	Type string `json:"type"`
	Expiration int64 `json:"exp"`
	JTI string `json:"jti"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(payload Claims, expiration time.Duration) (string, error) {
	jti := uuid.New().String()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload.Expiration = time.Now().Add(expiration).Unix()
	payload.JTI = jti

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

func VerifyJWT(w http.ResponseWriter, r *http.Request, isRefreshToken bool) (bool, *Claims) {
	var tokenString string
	if isRefreshToken {
		var requestBody struct {
			RefreshToken string `json:"refresh_token"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil || requestBody.RefreshToken == "" {
			http.Error(w, "Invalid request: missing refresh_token", http.StatusBadRequest)
			return false, nil
		}
		tokenString = requestBody.RefreshToken
	} else {
		tokenString = utils.GetJWTFromRequest(w, r)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg 
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing tokens
		return jwtSecret, nil
	})
	claims, ok := token.Claims.(*Claims)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Invalid token signature", http.StatusUnauthorized)
			return false, nil
		}
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return false, nil
	} else if !ok && !token.Valid {
		log.Printf("Token claims invalid or token is not valid")
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return false, nil
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

	return true, claims
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
