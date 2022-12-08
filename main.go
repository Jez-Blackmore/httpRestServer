package main

import (
	"fmt"
	log "httpRestServer/logging"
	req "httpRestServer/requests"
	store "httpRestServer/store"
	"net/http"
)

type Handler interface {
	ServeHTTP(response http.ResponseWriter, request *http.Request)
}

func main() {

	log.SetupLoggers()

	log.InfoLogger.Println("Creating store")
	store.MainStoreMain = store.NewStoreMain()

	go store.MainStoreMain.Monitor()

	log.InfoLogger.Println("Setting up REST endpoints")
	http.HandleFunc("/ping", req.ServerIsRunningGet)
	//http.HandleFunc("/login", req.Login)
	http.HandleFunc("/store/", req.UpdateStore)
	http.HandleFunc("/list", req.StoreList)
	http.HandleFunc("/list/", req.StoreListKey)
	http.HandleFunc("/shutdown", req.Shutdown)

	log.InfoLogger.Println("Starting Server")
	fmt.Println("Server Available - see http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.ErrorLogger.Println("Error starting server - ", err)
	}

}
