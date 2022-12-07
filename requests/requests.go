package requests

import (
	"encoding/json"
	log "httpRestServer/logging"
	"httpRestServer/store"
	"net/http"
	"os"
	"strings"
)

// GET /ping
// Public

func ServerIsRunningGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
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

	if key != "" {
		switch r.Method {
		case http.MethodGet:

			valueToShow := store.UpdateStoreGet(key)

			if valueToShow != "" {
				jsonData, err := json.Marshal(valueToShow)

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"Error"}`))
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(jsonData)

			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "404 key not found"}`))
			}

		case http.MethodPut:
			w.WriteHeader(http.StatusOK)

			var valueToUpdate store.StructValueObject
			valueToUpdate.Owner = username

			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&valueToUpdate); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message": "invalid payload"}`))
				return
			}

			value := store.UpdateStorePut(key, valueToUpdate)

			jsonData, err := json.Marshal(value)

			if err != nil {
				http.Error(w, "Error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)

		case http.MethodDelete:

			owner := store.GetKeyValueOwner(key)

			if owner != "" && username == owner {

				delete(store.Store, key)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return

			} else if username != owner {
				w.WriteHeader(http.StatusForbidden) // 403
				w.Write([]byte("Forbidden"))
				return
			}

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 key not found"))

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

		items := store.StoreListGet()

		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(items)

		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

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

			keyToShow, valueToShow := store.StoreListKeyGet(key)

			if keyToShow == "" {
				w.WriteHeader(http.StatusNotFound) //404
				w.Write([]byte(`404 key not found`))
				return
			} else {
				jsonData, err := json.Marshal(map[string]string{"key": keyToShow, "owner": valueToShow})

				if err != nil {
					http.Error(w, "Error", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(jsonData)
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
