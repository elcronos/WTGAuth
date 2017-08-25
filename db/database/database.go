package database

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"os"
	. "../data"
	. "../dataobjects"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DB	*gorm.DB
)

func init(){
	var err error

	if DB, err = OpenConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to database, but got err=%+v", err))
	}
	//Migrate the schema
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Country{})
	DB.AutoMigrate(&City{})
	//Insert Data after AutoMigrate
	InitialiseDB(DB)
}

func OpenConnection() (db *gorm.DB, err error){
	fmt.Println("Testing postgres...")
	db, err = gorm.Open("postgres", "client_encoding=UTF8 host=localhost user=postgres dbname=postgres sslmode=disable port=32768")

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}

	return db, err
}