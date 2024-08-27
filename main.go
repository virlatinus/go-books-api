package main

import (
	"encoding/json"
	"fmt"
	"github.com/MadAppGang/httplog"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	gorm.Model
	Title  string `json:"title" gorm:"not null;uniqueIndex:title_idx"`
	Author string `json:"author" gorm:"not null"`
	Year   int    `json:"year" gorm:"not null"`
}

type StatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PaginationResponse struct {
	TotalRecords int  `json:"total_records"`
	CurrentPage  int  `json:"current_page"`
	TotalPages   int  `json:"total_pages"`
	NextPage     *int `json:"next_page"`
	PrevPage     *int `json:"prev_page"`
	PageSize     int  `json:"page_size"`
}

type Response struct {
	Status     StatusResponse      `json:"status"`
	Data       interface{}         `json:"data"`
	Errors     []string            `json:"errors"`
	Pagination *PaginationResponse `json:"pagination"`
}

var db *gorm.DB

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
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256, // default size for string fields
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Book{})
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Use(httplog.Logger)
	router.Use(commonMiddleware)
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/batch", addBooks).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func logError(w http.ResponseWriter, err error, statusCode int) bool {
	if err != nil {
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(Response{
			Status: StatusResponse{
				Code:    statusCode,
				Message: "error",
			},
			Errors: []string{err.Error()},
		})
		return true
	}
	return false
}

func sendResponse(w http.ResponseWriter, statusCode int, response Response) bool {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(Response{
		Status: StatusResponse{
			Code:    statusCode,
			Message: "success",
		},
		Data:       response.Data,
		Pagination: response.Pagination,
	})
	return !logError(w, err, http.StatusInternalServerError)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	err := db.Find(&books).Error
	if logError(w, err, http.StatusInternalServerError) {
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Data: books,
		Pagination: &PaginationResponse{
			TotalRecords: len(books),
			CurrentPage:  1,
			TotalPages:   1,
			NextPage:     nil,
			PrevPage:     nil,
			PageSize:     len(books),
		},
	})
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	var book Book
	err = db.First(&book, id).Error
	if logError(w, err, http.StatusNotFound) {
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Data: book,
	})
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	err = db.Create(&book).Error
	if logError(w, err, http.StatusConflict) {
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Data: book,
	})
}

func addBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&books)
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	result := db.Create(&books)
	if logError(w, result.Error, http.StatusConflict) {
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Data: books,
	})
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	var book Book
	err = db.First(&book, id).Error
	if logError(w, err, http.StatusNotFound) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&book)
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	book.ID = uint(id) // preserve ID from params

	err = db.Save(&book).Error
	if logError(w, err, http.StatusConflict) {
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Data: book,
	})
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	var book Book
	err = db.First(&book, id).Error
	if logError(w, err, http.StatusNotFound) {
		return
	}

	err = db.Delete(&book).Error
	if logError(w, err, http.StatusInternalServerError) {
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Data: book,
	})
}
