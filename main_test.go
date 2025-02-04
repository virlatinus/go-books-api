package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_addBook(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{"My Book", args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer([]byte(`{"ID": 10, "Year": 1990, "Author": "John Doe", "Title": "My book"}`)))}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBook(tt.args.w, tt.args.r)
			if len(books) != 1 {
				t.Errorf("len(addBook()) = %d, want 1", len(books))
			}
		})
	}
}

func Test_deleteBook(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	tests := []struct {
		name string
		args args
	}{
		{"My Book", args{httptest.NewRecorder(), mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/books", nil), map[string]string{
			"id": "5",
		})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBooks()
			deleteBook(tt.args.w, tt.args.r)
			if len(books) != 4 {
				t.Errorf("len(deleteBook()) = %d, want 4", len(books))
			}
		})
	}
}

func Test_getBook(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{"My Book", args{httptest.NewRecorder(), mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/books", nil), map[string]string{
			"id": "1",
		})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBooks()
			getBook(tt.args.w, tt.args.r)
			r := tt.args.w.Result()
			if r.StatusCode != http.StatusOK {
				t.Errorf("getBook() status Code = %d, want %d", r.StatusCode, http.StatusOK)
			}
			expected := "{\"id\":1,\"title\":\"1984\",\"author\":\"George Orwell\",\"year\":1949}\n"
			if tt.args.w.Body.String() != expected {
				t.Errorf("getBook() = %q want %q",
					tt.args.w.Body.String(), expected)
			}
		})
	}
}

func Test_getBooks(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getBooks(tt.args.w, tt.args.r)
		})
	}
}

func Test_updateBook(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateBook(tt.args.w, tt.args.r)
		})
	}
}
