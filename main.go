package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Book struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Genre    string  `json:"genre"`
	Price    float64 `json:"price"`
	AuthorID string  `json:"author_id"`
}

type Author struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var jwtkey = []byte("secret_key")
var tokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)
var (
	authors = []Author{
		{ID: "1", FirstName: "John", LastName: "Doe"},
		{ID: "2", FirstName: "Jane", LastName: "Smith"},
		{ID: "3", FirstName: "Emily", LastName: "Johnson"},
	}

	books = []Book{
		{ID: "1", Title: "The Great Adventure", Genre: "Adventure", Price: 19.99, AuthorID: "1"},
		{ID: "2", Title: "Mystery of the Old House", Genre: "Mystery", Price: 14.99, AuthorID: "2"},
		{ID: "3", Title: "Learning Go", Genre: "Programming", Price: 29.99, AuthorID: "3"},
	}
)

var authorList map[string]Author = map[string]Author{}
var bookList map[string]Book = map[string]Book{}

func Init() {

	for _, val := range authors {
		authorList[val.ID] = val
	}

	for _, val := range books {
		bookList[val.ID] = val
	}

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	for _, val := range bookList {
		books = append(books, val)
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getBook(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	book, ok := bookList[params]
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
func getAuthors(w http.ResponseWriter, r *http.Request) {
	var authors []Author
	for _, val := range authorList {
		authors = append(authors, val)
	}
	err := json.NewEncoder(w).Encode(authors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAuthor(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	author, ok := authorList[params]
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

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
	}
	_, ok := bookList[book.ID]
	if ok {
		http.Error(w, "Book already exists", http.StatusConflict)
		return
	}
	bookList[book.ID] = book
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Book can not be added", http.StatusInternalServerError)
	}
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	_, ok := bookList[params]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}

	if book.ID != params {
		http.Error(w, "Inconsistent update request", http.StatusBadRequest)
		return
	}
	bookList[params] = book
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Book can not be updated", http.StatusInternalServerError)
		return
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	_, ok := bookList[params]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	delete(bookList, params)
	_, err := w.Write([]byte("Book deleted"))
	if err != nil {
		http.Error(w, "Book can not be deleted", http.StatusInternalServerError)
	}

}

func getTopBooks(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "limit")
	limit, err := strconv.Atoi(params)
	if err != nil {
		http.Error(w, "error while parsing limit", http.StatusBadRequest)
	}
	var books []Book
	for _, val := range bookList {
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

func logIn(w http.ResponseWriter, r *http.Request) {
	var cred Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}
	expectedPassword, ok := users[cred.Username]
	fmt.Println(expectedPassword, ok)

	if !ok || expectedPassword != cred.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(time.Hour * 2)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"aud": cred.Username,
		"exp": expirationTime.Unix(),
	})
	if err != nil {
		http.Error(w, "Can not generate jwt", http.StatusInternalServerError)
		return
	}
	fmt.Println(tokenString)

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Successfully Logged In"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func logOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Successfully Logged Out"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	Init()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/ping", ping)
	r.Post("/login", logIn)
	r.Post("/logout", logOut)

	r.Group(func(r chi.Router) {
		r.Route("/book", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Post("/", addBook)
			r.Put("/{id}", updateBook)
			r.Delete("/{id}", deleteBook)
		})
	})

	r.Get("/books", getBooks)
	r.Get("/book/{id}", getBook)
	r.Get("/authors", getAuthors)
	r.Get("/author/{id}", getAuthor)
	r.Get("/books/top/{limit}", getTopBooks)
	fmt.Println("Listening and Serving to 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
		return
	}

}
