package api

import (
	"net/http"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/api/middleware"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/handler"
)

type Middleware func(http.Handler) http.Handler

func Routes() *http.ServeMux {
	router := http.NewServeMux()

	// send a POST request with a form data of email and password to this endpoint
	// using Content-Type: application/x-www-form-urlencoded or multipart/form-data.

	router.HandleFunc("/auth/sign-up", handler.SignUpHandler)
	router.HandleFunc("/auth/refresh", handler.RefreshHandler)
	router.HandleFunc("/auth/sign-in", handler.SignInHandler)
	router.Handle("/service", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You are authenticated"))
	})))
	return router
}

func middleWareChain(h http.Handler, middlewares ...Middleware) http.Handler {
	// Iterate backwards
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
