package main

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/pbkdf2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseProblem struct {
	Msg string `json:"msg"`
}

func InitialMigration() {
	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	//Migrate database
	db.AutoMigrate(&User{})
}

func UserRegiter(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	var person_struct User

	err = json.NewDecoder(r.Body).Decode(&person_struct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(person_struct.Username) == 0 || len(person_struct.Email) == 0 || len(person_struct.Password) == 0 {
		response_problem := ResponseProblem{"Error with json in the body!"}

		js_response, err_json := json.Marshal(response_problem)
		if err_json != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		w.Write(js_response)
	} else {
		// Check if the desired user exists:
		if user_err := db.Where("username = ?", person_struct.Username).First(&person_struct).Error; user_err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			salt := []byte("mysaltedSalt")
			derivedKey := pbkdf2.Key([]byte(person_struct.Password), salt, 10, 256, sha256.New)
			person_struct.Password = string(derivedKey)
			db.Create(&person_struct)
			db.Save(&person_struct)
			js_response, err_json := json.Marshal(person_struct)
			if err_json != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Write(js_response)

		} else if user_err == nil {
			response_problem := ResponseProblem{"User alredy exists!"}
			js_response, err_json := json.Marshal(response_problem)
			if err_json != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(js_response)
		}
	}
}
