package main

import "github.com/gorilla/mux"

func main() {
	//Migrate db:
	InitialMigration()
	//MUX Router:
	router := mux.NewRouter()

}
