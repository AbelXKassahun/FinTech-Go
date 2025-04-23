package middleware

import (
	"log"
	"net/http"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Println("in auth middleware 2")
		ok, _ := auth.VerifyJWT(w, r, false)
		if !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}
