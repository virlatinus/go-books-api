package controllers

import (
	. "books-api/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type NewController struct {
	DB *gorm.DB
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

func (bp *NewController) GetBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	err := bp.DB.Find(&books).Error
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

func (bp *NewController) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	var book Book
	err = bp.DB.First(&book, id).Error
	if logError(w, err, http.StatusNotFound) {
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Data: book,
	})
}

func (bp *NewController) AddBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	err = bp.DB.Create(&book).Error
	if logError(w, err, http.StatusConflict) {
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Data: book,
	})
}

func (bp *NewController) AddBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&books)
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	result := bp.DB.Create(&books)
	if logError(w, result.Error, http.StatusConflict) {
		return
	}

	sendResponse(w, http.StatusCreated, Response{
		Data: books,
	})
}

func (bp *NewController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	var book Book
	err = bp.DB.First(&book, id).Error
	if logError(w, err, http.StatusNotFound) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&book)
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	book.ID = uint(id) // preserve ID from params

	err = bp.DB.Save(&book).Error
	if logError(w, err, http.StatusConflict) {
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Data: book,
	})
}

func (bp *NewController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if logError(w, err, http.StatusBadRequest) {
		return
	}

	var book Book
	err = bp.DB.First(&book, id).Error
	if logError(w, err, http.StatusNotFound) {
		return
	}

	err = bp.DB.Delete(&book).Error
	if logError(w, err, http.StatusInternalServerError) {
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Data: book,
	})
}
