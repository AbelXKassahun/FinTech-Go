package api

import (
	"net/http"
	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/handler"
)

func Routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/auth/sign-up", handler.SignUpHandler)
	// router.HandleFunc("/auth/sign-in", SignInHandler)

	return router
}