package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Username string
	Email    string
	Password string
}

func InitialMigration() {
	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	//Migrate database
	db.AutoMigrate(&User{})
}
