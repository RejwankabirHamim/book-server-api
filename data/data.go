package data

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
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

var Users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var jwtkey = []byte("secret_key")
var TokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)
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

var AuthorList map[string]Author = map[string]Author{}
var BookList map[string]Book = map[string]Book{}

func Init() {

	for _, val := range authors {
		AuthorList[val.ID] = val
	}

	for _, val := range books {
		BookList[val.ID] = val
	}

}
