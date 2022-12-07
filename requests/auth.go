package requests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("my_secret_key_123_a_bit_better") // only suitable for dev

// GET /login
// Private - Authorisation required

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "charset=utf-8")

	username, password, ok := r.BasicAuth()

	if !ok || username == "" || password == "" {
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	switch r.Method {
	case http.MethodGet:

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			fmt.Printf("error %v", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "not found"}`))
			return
		}

		fmt.Println("token2: ", tokenString, err)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Bearer":` + tokenString + `}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}
