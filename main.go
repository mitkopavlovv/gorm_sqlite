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
	port := "8000"
	fmt.Println("Server started on port " + port)
	log.Println(http.ListenAndServe(":"+port, router))

}
