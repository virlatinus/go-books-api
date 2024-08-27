package main

import (
	"encoding/json"
	"github.com/MadAppGang/httplog"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"slices"
	"strconv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

var books []Book

func main() {
	router := mux.NewRouter()

	books = append(books,
		Book{ID: 1, Title: "1984", Author: "George Orwell", Year: 1949},
		Book{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee", Year: 1960},
		Book{ID: 3, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Year: 1925},
		Book{ID: 4, Title: "Pride and Prejudice", Author: "Jane Austen", Year: 1813},
		Book{ID: 5, Title: "The Catcher in the Rye", Author: "J.D. Salinger", Year: 1951},
	)

	router.Use(httplog.Logger)
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slices.SortFunc(books, func(a, b Book) int { return a.ID - b.ID })
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	idx := slices.IndexFunc(books, func(b Book) bool { return b.ID == id })
	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books[idx])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	idx := slices.IndexFunc(books, func(b Book) bool { return b.ID == book.ID })
	if idx != -1 {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte("Book ID already exists\n"))
		return
	}

	books = append(books, book)
	w.WriteHeader(http.StatusCreated)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	idx := slices.IndexFunc(books, func(b Book) bool { return b.ID == id })
	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var book Book
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	book.ID = id // keep the same ID
	books[idx] = book
	w.WriteHeader(http.StatusNoContent)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	idx := slices.IndexFunc(books, func(b Book) bool { return b.ID == id })
	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	books = slices.Delete(books, idx, idx+1)
	w.WriteHeader(http.StatusNoContent)
}
