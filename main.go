package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var jwtkey = []byte("secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

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

func Login(w http.ResponseWriter, r *http.Request) {
	var credential Credential
	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword, ok := users[credential.Username]
	if !ok || expectedPassword != credential.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 4)
	claims := &Claims{
		Username: credential.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenString := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("Hello %s!", claims.Username)))
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

func main() {
	Init()
	r := chi.NewRouter()
	r.Get("/books", getBooks)
	r.Get("/book/{id}", getBook)
	r.Get("/authors", getAuthors)
	r.Get("/author/{id}", getAuthor)
	r.Post("/book", addBook)
	r.Put("/book/{id}", updateBook)
	r.Delete("/book/{id}", deleteBook)
	r.Get("/books/top/{limit}", getTopBooks)
	r.Post("/login", Login)
	r.Get("/home", Home)
	fmt.Println("Listening and Serving to 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
