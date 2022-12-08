package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
)

type Users struct {
	Users []User
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

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

		CheckUserMatches()

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

func CheckUserMatches() {

	abs, err := filepath.Abs("./requests/users.json")

	jsonFile, err := os.Open(abs)
	fmt.Println("Successfully Opened users.json")

	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users Users

	err = json.Unmarshal(byteValue, &users)

	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(users.Users); i++ {
		fmt.Println("User Type: " + users.Users[i].Username)
		fmt.Println("User Name: " + users.Users[i].Password)
		fmt.Println("Facebook Url: " + users.Users[i].Role)
	}

}
