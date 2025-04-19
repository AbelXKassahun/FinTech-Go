package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/auth"
	"github.com/google/uuid"
)

// [To Do]
// check if the user exists in SignUpHandler and SignInHandler
// refactor the form validation and token response
// implement the logout handler
// store user metadata (e.g., tier, last login, session history) in Redis hash structure.
// get the user ip and device name

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
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	}
	if !auth.IsValidPassword(password) || password == "" {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	// [To Do] - check if the user exists

	userID := uuid.New().String()
	// generate the access token
	accessToken, _, err := auth.GenerateJWT(auth.Claims{
		UserID: userID,
		Tier:  "basic",
		Type:  "access_token",
	}, 15*time.Minute)

	if err != nil {
		http.Error(w, "Error generating access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// generate the refresh token
	refreshToken, _, err := auth.GenerateJWT(auth.Claims{
		UserID: userID,
		Tier:  "basic",
		Type:  "refresh_token",
	}, 7*24*time.Hour)

	if err != nil {
		http.Error(w, "Error generating refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		http.Error(w, "Couldnt encode tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/auth/sign-in called")

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
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	}
	if !auth.IsValidPassword(password) || password == "" {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	// [To Do] - check if the user exists and validate the password here

	userID := uuid.New().String()

	accessToken, _, err := auth.GenerateJWT(auth.Claims{
		UserID: userID,
		Tier:  "basic",
		Type:  "access_token",
	}, 15*time.Minute)

	if err != nil {
		http.Error(w, "Error generating access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, _,  err := auth.GenerateJWT(auth.Claims{
		UserID: userID,
		Tier:  "basic",
		Type:  "refresh_token",
	}, 7*24*time.Hour)

	if err != nil {
		http.Error(w, "Error generating refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		http.Error(w, "Couldnt encode tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {

}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/auth/refresh called")

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	ok, claims := auth.VerifyJWT(w, r, true)
	if !ok {
		return
	}
	log.Printf("Refresh token validated for user '%s'", claims.UserID)


	newAccessToken, _, err := auth.GenerateJWT(*claims, 15*time.Minute)
	if err != nil {
		log.Printf("Error generating new access tokens during refresh: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	newRefreshToken, _, err := auth.GenerateJWT(*claims, 7*24*time.Hour)
	if err != nil {
		log.Printf("Error generating new refresh tokens during refresh: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("Refresh and access token generated for user '%s'", claims.UserID)
	// 5. Send new tokens back
	tokens := TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken, // Implement refresh token rotation
		TokenType:    "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		log.Println("Couldnt encode tokens: " + err.Error())
		http.Error(w, "Couldnt encode tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}
}