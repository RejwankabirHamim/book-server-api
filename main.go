package main

import (
	"fmt"
	"github.com/RejwankabirHamim/api-book-server/auth"
	"github.com/RejwankabirHamim/api-book-server/data"
	"github.com/RejwankabirHamim/api-book-server/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
)

func main() {
	data.Init()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/ping", handler.Ping)
	r.Post("/login", auth.LogIn)
	r.Post("/logout", auth.LogOut)

	r.Group(func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", handler.GetBooks)
			r.Get("/{id}", handler.GetBook)
			r.Get("/top/{limit}", handler.GetTopBooks)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(data.TokenAuth))
				r.Use(jwtauth.Authenticator(data.TokenAuth))

				r.Post("/", handler.AddBook)
				r.Put("/{id}", handler.UpdateBook)
				r.Delete("/{id}", handler.DeleteBook)
			})

		})

		r.Route("/authors", func(r chi.Router) {
			r.Get("/", handler.GetAuthors)
			r.Get("/{id}", handler.GetAuthor)
		})
	})

	fmt.Println("Listening and Serving to 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
		return
	}

}
