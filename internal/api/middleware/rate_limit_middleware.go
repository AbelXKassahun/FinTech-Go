package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/auth"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetJWTFromRequest(w, r)
		claims, err := getClaimsFromTokenString(tokenString)
		jti := claims["ID"].(string)
		if err != nil {
			log.Println("couldnt get claims from token",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		proceed, err := auth.CheckToken(r.Context(), jti)
		if err != nil {
			log.Println("couldnt check token",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !proceed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LogInRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		proceed, err := auth.CheckToken(r.Context(), ip)

		if err != nil {
			log.Println("couldnt check token",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !proceed {
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
