package handler

import (
	"encoding/json"
	"github.com/RejwankabirHamim/api-book-server/data"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sort"
	"strconv"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	var books []data.Book
	for _, val := range data.BookList {
		books = append(books, val)
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func GetBook(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	book, ok := data.BookList[params]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func GetAuthors(w http.ResponseWriter, r *http.Request) {
	var authors []data.Author
	for _, val := range data.AuthorList {
		authors = append(authors, val)
	}
	err := json.NewEncoder(w).Encode(authors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAuthor(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	author, ok := data.AuthorList[params]
	if !ok {
		http.Error(w, "Author not found", http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	var book data.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
	}
	_, ok := data.BookList[book.ID]
	if ok {
		http.Error(w, "Book already exists", http.StatusConflict)
		return
	}
	data.BookList[book.ID] = book
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Book can not be added", http.StatusInternalServerError)
	}
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	_, ok := data.BookList[params]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	var book data.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}

	if book.ID != params {
		http.Error(w, "Inconsistent update request", http.StatusBadRequest)
		return
	}
	data.BookList[params] = book
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Book can not be updated", http.StatusInternalServerError)
		return
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	_, ok := data.BookList[params]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	delete(data.BookList, params)
	_, err := w.Write([]byte("Book deleted"))
	if err != nil {
		http.Error(w, "Book can not be deleted", http.StatusInternalServerError)
	}

}

func GetTopBooks(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "limit")
	limit, err := strconv.Atoi(params)
	if err != nil {
		http.Error(w, "error while parsing limit", http.StatusBadRequest)
	}
	var books []data.Book
	for _, val := range data.BookList {
		books = append(books, val)
	}
	sort.Slice(books, func(i, j int) bool {
		return books[i].Price > books[j].Price
	})
	top := books[:limit]
	err = json.NewEncoder(w).Encode(top)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
