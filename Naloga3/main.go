package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"netis_naloga/httpHandler"
)

func main() {
	router := mux.NewRouter()
	// HANDLE ENDPOINTS
	router.HandleFunc("/sign", httpHandler.Sign).Methods("POST")
	router.HandleFunc("/public", httpHandler.Public).Methods("POST")
	router.HandleFunc("/validate", httpHandler.Validate).Methods("POST")
	// RUN SERVER
	err := http.ListenAndServe(":8080", router)
	fmt.Println("Server is listening on  port 8080.")
	if err != nil {
		panic(err)
	}
}
