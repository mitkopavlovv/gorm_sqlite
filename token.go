package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/pbkdf2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TokenMessage struct {
	Token string `json:"token"`
}

var signingKey = []byte("S00perS3cret")

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = username
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		fmt.Errorf("JWT Generation failed: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ReturnToken(w http.ResponseWriter, r *http.Request) {
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

	// Check if user and pass exists exists in database:
	salt := []byte("mysaltedSalt")
	derivedKey := pbkdf2.Key([]byte(person_struct.Password), salt, 10, 256, sha256.New)
	if user_err := db.Where("username = ? AND password = ?", person_struct.Username, string(derivedKey)).First(&person_struct).Error; user_err == nil {
		token, err := GenerateToken(person_struct.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonToken := TokenMessage{token}
		js_response, err_json := json.Marshal(jsonToken)
		if err_json != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js_response)

	}
}

func AuthMiddleware(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error with signing key")
				}
				return signingKey, nil
			})
			if err != nil {
				response_problem := ResponseProblem{"Error with signing key"}
				js_response, err_json := json.Marshal(response_problem)
				if err_json != nil {
					http.Error(w, err_json.Error(), http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(js_response)
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			response_problem := ResponseProblem{"Not authorized"}
			js_response, err_json := json.Marshal(response_problem)
			if err_json != nil {
				http.Error(w, err_json.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(js_response)
		}
	})
}
