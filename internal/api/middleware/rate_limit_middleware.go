package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/auth"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/storage"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetJWTFromRequest(w, r)
		claims, err := getClaimsFromTokenString(tokenString)
		if err != nil {
			log.Println("couldnt get claims from token", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jti := claims["jti"].(string)
		log.Println("jti", jti)
		proceed, err := auth.CheckToken(r.Context(), jti)
		if err != nil && err != storage.RedisErr {
			log.Println("couldnt check token", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err == storage.RedisErr { // first time creating rate token for user
			err = auth.NewRateToken(r.Context(), jti, time.Minute*15, float64(3), 0.5)
			if err != nil {
				http.Error(w, "Error generating rate token: "+err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		if !proceed && err == nil {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LogInRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("in rate limit middleware 2")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("Couldn't get ip from request", err)
			ip = r.RemoteAddr
		}

		proceed, err := auth.CheckToken(r.Context(), ip)

		if err != nil && err != storage.RedisErr {
			log.Println("Couldn't check token", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err == storage.RedisErr { // first login request
			err = auth.NewRateToken(r.Context(), ip, time.Minute*15, float64(3), 0.5)
			if err != nil {
				http.Error(w, "Error generating rate token: "+err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// subsequent login requests
		if !proceed && err == nil {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getClaimsFromTokenString(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
