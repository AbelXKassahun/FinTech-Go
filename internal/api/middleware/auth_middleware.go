package middleware

import (
	"net/http"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		ok, _ := auth.VerifyJWT(w, r, false)
		if !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}
