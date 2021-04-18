package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Msg string `json:"msg"`
}

func SecretData(w http.ResponseWriter, r *http.Request) {
	resp := Response{"Super Sectet Message"}
	json_resp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_resp)
}

func main() {
	//Migrate db:
	InitialMigration()
	//MUX Router:
	router := mux.NewRouter()
	router.HandleFunc("/api/register", UserRegiter).Methods("POST")
	router.HandleFunc("/api/token", ReturnToken).Methods("POST")
	router.Handle("/api/secret", AuthMiddleware(SecretData)).Methods("GET")
	port := "8000"
	fmt.Println("Server started on port " + port)
	log.Println(http.ListenAndServe(":"+port, router))
}
