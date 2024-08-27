package main

import (
	"books-api/controllers"
	"books-api/middleware"
	"books-api/models"
	"fmt"
	"github.com/MadAppGang/httplog"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	var err error

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.Get("DB_USER"),
		viper.Get("DB_PASS"),
		viper.Get("DB_HOST"),
		viper.Get("DB_PORT"),
		viper.Get("DB_NAME"),
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256, // default size for string fields
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Use(httplog.Logger)
	router.Use(middleware.JsonHeader)

	bc := controllers.NewController{DB: db}
	router.HandleFunc("/books", bc.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", bc.GetBook).Methods("GET")
	router.HandleFunc("/books", bc.AddBook).Methods("POST")
	router.HandleFunc("/books/batch", bc.AddBooks).Methods("POST")
	router.HandleFunc("/books/{id}", bc.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", bc.DeleteBook).Methods("DELETE")

	log.Printf("Starting server on %s\n", viper.Get("APP_ADDRESS"))
	log.Fatal(http.ListenAndServe(viper.GetString("APP_ADDRESS"), router))
}
