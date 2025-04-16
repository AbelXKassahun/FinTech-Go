package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/auth"
)

// [To Do] - check if the user exists in SignUpHandler

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	if !auth.IsValidEmail(email) || email == "" {
		// w.WriteHeader(400)
		http.Error(w, "Invalid Email", http.StatusBadRequest)
	}
	if !auth.IsValidPassword(password) || password == "" {
		// w.WriteHeader(400)
		http.Error(w, "Invalid Password", http.StatusBadRequest)
	}

	// [To Do] - check if the user exists

	// generate the access token
	accessToken, err := auth.GenerateJWT(map[string]interface{}{
		"email": email,
		"tier":  "basic",
		"typ":   "access_token",
	}, 15*time.Minute)

	if err != nil {
		http.Error(w, "Error generating access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// generate the refresh token
	refreshToken, err := auth.GenerateJWT(map[string]interface{}{
		"email": email,
		"tier":  "basic",
		"typ":   "refresh_token",
	}, 7*24*time.Hour)

	if err != nil {
		http.Error(w, "Error generating refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType: "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		http.Error(w, "Couldnt encode tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
