package main

import (
	"fmt"
	log "httpRestServer/logging"
	req "httpRestServer/requests"
	"net/http"
)

type Handler interface {
	ServeHTTP(response http.ResponseWriter, request *http.Request)
}

func main() {

	log.SetupLoggers()

	log.InfoLogger.Println("Starting Server")

	fmt.Printf("%s", req.Store["test"])

	http.HandleFunc("/ping", req.ServerIsRunningGet)
	http.HandleFunc("/login", req.Login)
	http.HandleFunc("/store/", req.UpdateStore)
	http.HandleFunc("/list", req.StoreList)
	http.HandleFunc("/list/", req.StoreListKey)
	http.HandleFunc("/shutdown", req.Shutdown)

	fmt.Println("Server Available - see http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.ErrorLogger.Println("Error starting server - ", err)
	}

}
