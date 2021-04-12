package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//Migrate db:
	InitialMigration()
	//MUX Router:
	router := mux.NewRouter()
	router.HandleFunc("/api/register", UserRegiter).Methods("POST")
	fmt.Println("Server started")
	log.Println(http.ListenAndServe(":8000", router))

}
