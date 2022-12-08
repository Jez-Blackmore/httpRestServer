package requests

import (
	"encoding/json"
	"fmt"
	log "httpRestServer/logging"
	"httpRestServer/store"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// GET /ping
// Public

func ServerIsRunningGet(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
		return
	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// GET/PUT/DELETE /store/<key>
// Private - Authorisation required

func UpdateStore(w http.ResponseWriter, r *http.Request) {

	fmt.Print("RUNNNING Again")
	username := r.Header.Get("Authorization")

	if username == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/store/")

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}

	responseString := string(responseData)

	if key != "" {
		switch r.Method {
		case http.MethodGet:

			valueToShow := store.MainStoreMain.UpdateStoreGet(key)

			if valueToShow != "" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(valueToShow))
				return

			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 key not found"))
				return
			}

		case http.MethodPut:
			var valueToUpdate store.StructValueObject
			valueToUpdate.Owner = username

			valueToUpdate.Value = responseString

			value := store.MainStoreMain.UpdateStorePut(store.Key(key), valueToUpdate, username)

			if value != "" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(value))
				return
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8") //403
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Forbidden"))
				return
			}

		case http.MethodDelete:

			owner := store.MainStoreMain.GetKeyValueOwner(key)

			if owner == "" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 key not found"))
				return
			}

			if username == owner || username == "admin" {

				success := store.MainStoreMain.UpdateStoreDelete(store.Key(key))

				if success {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusOK)
					return
				} else {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusForbidden) // 403
					return
				}

			} else if username != owner {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusForbidden) // 403
				return
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 key not found"))
				return
			}

		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 key not found"))
			return
		}
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("You must provide a key"))
		return
	}
}

// GET /list - should return a list
// Private - Authorisation required

func StoreList(w http.ResponseWriter, r *http.Request) {

	username := r.Header.Get("Authorization")

	if username == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	switch r.Method {
	case http.MethodGet:

		items := store.MainStoreMain.StoreListGet()

		jsonData, err := json.Marshal(items)

		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
		return
	}
}

// GET /list/<key> - should return infomration based on list key
// Private - Authorisation required

func StoreListKey(w http.ResponseWriter, r *http.Request) {

	key := strings.TrimPrefix(r.URL.Path, "/list/")

	username := r.Header.Get("Authorization")

	if username == "" {
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	if key != "" {
		switch r.Method {
		case http.MethodGet:

			keyToShow, valueToShow := store.MainStoreMain.StoreListKeyGet(key)

			if keyToShow == "" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusNotFound) //404
				w.Write([]byte(`404 key not found`))
				return
			} else {
				jsonData, err := json.Marshal(map[string]string{"key": keyToShow, "owner": valueToShow})

				if err != nil {
					http.Error(w, "Error", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonData)
				return
			}

		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 key not found")) // 404
			return
		}
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 key not found")) // 404
		return
	}
}

// GET /shutdown
// Private - Authorisation required

func Shutdown(w http.ResponseWriter, r *http.Request) {

	username := r.Header.Get("Authorization")

	if username == "" || username != "admin" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		log.InfoLogger.Println("Shutting Down Server")
		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(0)
		}()

	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
		return
	}
}
