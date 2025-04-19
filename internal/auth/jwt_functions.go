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
	UserID     string `json:"user_id"`
	Tier       string `json:"tier"`
	Type       string `json:"type"`
	jwt.RegisteredClaims
}

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(payload Claims, expiration time.Duration) (string, *Claims, error) {
	jti := uuid.New().String()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}


	payload.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		ID:        jti,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't marshal header: %v", err)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't marshal payload: %v", err)
	}

	headerEncoded := utils.Base64URLEncode(headerJSON)
	payloadEncoded := utils.Base64URLEncode(payloadJSON)

	unsigned := headerEncoded + "." + payloadEncoded
	signature := utils.Sign(unsigned, JWTSecret)

	return unsigned + "." + signature, &payload, nil
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

	// verify signiture and parse claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing tokens
		return JWTSecret, nil
	})
	
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Invalid token signature", http.StatusUnauthorized)
			return false, nil
		}
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return false, nil
	}

	// verifying standard claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Printf("Token claims invalid or token is not valid")
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return false, nil
	}

	// [To Do] - uncomment this code if you want to check if the token is blacklisted

	jti := claims.ID
	log.Println(jti)

	blacklisted, err := IsBlacklisted(jti, r.Context())
	if blacklisted {
		w.WriteHeader(401)
		w.Write([]byte("Token is invalid"))
		return false, nil
	} else if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return false, nil
	}

	return true, claims
}

// when the user logs out or changes password
func BlacklistToken(tokenString string, ctx context.Context) error {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	claims := token.Claims.(jwt.MapClaims)
	jti := claims["jti"].(string)
	// get the remaining time until the token expires
	exp := time.Until(time.Unix(int64(claims["exp"].(float64)), 0))

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
