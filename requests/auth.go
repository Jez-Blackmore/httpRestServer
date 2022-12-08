package requests

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Users struct {
	Users []User
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NEVER do this. Like ever. Argon2 or Bcrypt the passwords.
var users = map[string]string{
	"user_a": "passwordA",
	"user_b": "passwordB",
	"user_c": "passwordC",
	"admin":  "Password1",
}

var jwtKey = []byte("my_secret_key_123_a_bit_better") // only suitable for dev

// GET /login
// Private - Authorisation required

func Login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	if !ok {
		w.Header().Set("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusForbidden) // 401
		w.Write([]byte(`Forbidden`))
		return
	}

	if username == "" || password == "" {
		w.Header().Set("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusForbidden) // 401
		w.Write([]byte(`Forbidden`))
		return
	}
	/* bearerToken := r.Header.Get("Authorization")
	token := strings.TrimPrefix(bearerToken, "Basic ") */

	/* fmt.Println("4444444: ", token) */

	switch r.Method {
	case http.MethodGet:

		valid := CheckUserMatches(username, password)

		if !valid {
			w.Header().Set("Content-Type", "charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`Unauthorized`))
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"password": password,
			/* 	"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(), */
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			fmt.Printf("error %v", err)
			w.Header().Set("Content-Type", "charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "not found"}`))
			return
		}

		w.Header().Set("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Bearer ` + tokenString))
		return
	default:
		w.Header().Set("Content-Type", "charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`Unauthorized`))
		return
	}
}

func CheckUserMatches(username string, password string) bool {

	for key, value := range users {
		if key == username && password == value {
			/* 	fmt.Println("user: ", key, value) */
			return true
		}
	}

	return false

}

func ValidateUser(basicAuthVal string) (bool, string) {

	tokenString := strings.TrimPrefix(basicAuthVal, "Bearer ")

	/* fmt.Println("this: ", tokenString) */

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	var usernameToCheck string
	var posswordToCheck string

	// do something with decoded claims
	for key, val := range claims {
		if key == "password" {
			posswordToCheck = val.(string)
		}
		if key == "username" {
			usernameToCheck = val.(string)
		}
		/* fmt.Printf("Key: %v, value: %v\n", key, val) */
	}

	isValuseUser := CheckUserMatches(usernameToCheck, posswordToCheck)

	return isValuseUser, usernameToCheck
}
