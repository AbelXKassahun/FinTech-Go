package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"net/http"
)

func Base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func Sign(message string, jwtSecret []byte) string {
	h := hmac.New(sha256.New, jwtSecret)
	h.Write([]byte(message))
	return Base64URLEncode(h.Sum(nil))
}

func GetJWTFromRequest(w http.ResponseWriter, r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(498)
		w.Write([]byte("authorization header missing"))
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		w.WriteHeader(498)
		w.Write([]byte("authorization header format must be Bearer {token}"))
		return ""
	}
	token := strings.TrimPrefix(authHeader, prefix)
	return token
}