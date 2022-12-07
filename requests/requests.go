package requests

import (
	"encoding/json"
	"fmt"
	log "httpRestServer/logging"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("my_secret_key_123_a_bit_better") // only suitable for dev

var (
	Store = map[string]string{"Test": "test"}
)

type StructValueObject struct {
	Value string `json:"value"`
	Owner string `json:"owner"`
}

// GET /ping
// Public

func ServerIsRunningGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		fmt.Println("Get")
		w.Write([]byte("pong"))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// GET/PUT/DELETE /store/<key>
// Private - Authorisation required

func UpdateStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "charset=utf-8")

	username, _, ok := r.BasicAuth()

	if !ok {
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/store/")
	fmt.Println("key:", key)

	fmt.Println(Store)

	if key != "" {
		switch r.Method {
		case http.MethodGet:

			var valueToShow string

			for keyVal, value := range Store {
				if keyVal == key {
					valueToShow = value
				}
			}

			if valueToShow != "" {
				jsonData, err := json.Marshal(valueToShow)

				if err != nil {
					fmt.Println("Problem", err)
					http.Error(w, "Error",
						http.StatusInternalServerError)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message": "error"}`))
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write(jsonData)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "404 key not found"}`))
			}

		case http.MethodPut:
			w.WriteHeader(http.StatusOK)

			var keyToShow string

			for keyVal, _ := range Store {
				if keyVal == key {
					keyToShow = keyVal
				}
			}

			fmt.Println("username: ", username)

			var valueToUpdate StructValueObject
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&valueToUpdate); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "invalid payload"}`))
				return
			}

			if keyToShow == "" {
				Store[key] = valueToUpdate.Value
			} else {
				Store[keyToShow] = valueToUpdate.Value
			}

			jsonData, err := json.Marshal(valueToUpdate.Value)

			if err != nil {
				fmt.Println("Problem", err)
				http.Error(w, "Error",
					http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(jsonData)
			}

		case http.MethodDelete:

			keyExists := doesKeyExistInStore(key)

			if keyExists && username == "jez" {

				delete(Store, key)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return

			} else if username != "jez" {
				w.WriteHeader(http.StatusForbidden) // 403
				w.Write([]byte("Forbidden"))
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 key not found"))
				return
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "not found"}`))
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("You must provide a key"))
		return
	}
}

// GET /list - should return a list
// Private - Authorisation required

func StoreList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "charset=utf-8")

	username, _, ok := r.BasicAuth()

	if !ok || username == "" {
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	switch r.Method {
	case http.MethodGet:

		var i = 0
		for keyVal, value := range Store {
			fmt.Println(i, keyVal)
			fmt.Println(i, value)
			i++
		}

		w.WriteHeader(http.StatusOK)

		jsonData, err := json.Marshal(Store)

		if err != nil {
			fmt.Println("Problem", err)
			http.Error(w, "Error",
				http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
			return
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
		return
	}
}

// GET /list/<key> - should return infomration based on list key
// Private - Authorisation required

func StoreListKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "charset=utf-8")

	key := strings.TrimPrefix(r.URL.Path, "/list/")

	username, _, ok := r.BasicAuth()

	if !ok || username == "" {
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	if key != "" {
		switch r.Method {
		case http.MethodGet:

			var keyToShow string
			var valueToShow string

			for keyVal, value := range Store {
				if keyVal == key {
					keyToShow = keyVal
					valueToShow = value
					fmt.Printf("%s value is %v\n", keyVal, value)
				}
			}

			if keyToShow == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "not found"}`))
				return
			} else {
				jsonData, err := json.Marshal(map[string]string{"key": keyToShow, "owner": valueToShow})

				if err != nil {
					fmt.Println("Problem", err)
					http.Error(w, "Error",
						http.StatusInternalServerError)
					return
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write(jsonData)
					return
				}
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 key not found")) // 404
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 key not found")) // 404
		return
	}
}

// GET /shutdown
// Private - Authorisation required

func Shutdown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "charset=utf-8")

	username, _, ok := r.BasicAuth()

	if !ok || username == "" {
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	switch r.Method {
	case http.MethodGet:

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		log.InfoLogger.Println("Shutting Down Server")
		os.Exit(0)

	default:
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}
}

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

func doesKeyExistInStore(key string) bool {

	for keyVal, _ := range Store {
		if keyVal == key {
			return true
		}
	}

	return false
}
