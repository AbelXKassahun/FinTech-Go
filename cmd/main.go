package main

import (
	// "context"
	"fmt"
	// "log"
	// "net/http"

	// "github.com/AbelXKassahun/Digital-Wallet-Platform/internal/api"
)

var PORT = ":8080"
func main() {
	// log.Printf("Listening at port 8080 \n")
	// if err := http.ListenAndServe(PORT, api.Routes()); err != nil {
	// 	panic(err)
	// }

	ff := true
	key := "outside"
	
	if ff {
		key = "inside"
	}

	fmt.Println(key)
}


/*
// dummy 
	payload := map[string]interface{}{
		"tier": "basic",
	}
	expiration := time.Minute * 15
	// context := context.Background()
	// fmt.Println(auth.GenerateJWT(payload, expiration))
	// fmt.Println(auth.VerifyJWT("", context.Background()))

	router := http.NewServeMux()
	router.HandleFunc("GET /auth/get-jwt", func(w http.ResponseWriter, r *http.Request) {
		jwt_token, err := auth.GenerateJWT(payload, expiration)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Couldn't generate jwt token"))
			return
		}

		w.Write([]byte(jwt_token))
	})

	router.HandleFunc("GET /auth/verify-jwt", func(w http.ResponseWriter, r *http.Request) {
		auth.VerifyJWT(w, r)
	})

*/