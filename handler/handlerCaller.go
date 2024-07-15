package handler

import (
	"fmt"
	"github.com/RejwankabirHamim/api-book-server/auth"
	"github.com/RejwankabirHamim/api-book-server/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
)

func Caller(Port string) {
	data.Init()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/ping", Ping)
	r.Post("/login", auth.LogIn)
	r.Post("/logout", auth.LogOut)

	r.Group(func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", GetBooks)
			r.Get("/{id}", GetBook)
			r.Get("/top/{limit}", GetTopBooks)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(data.TokenAuth))
				r.Use(jwtauth.Authenticator(data.TokenAuth))

				r.Post("/", AddBook)
				r.Put("/{id}", UpdateBook)
				r.Delete("/{id}", DeleteBook)
			})

		})

		r.Route("/authors", func(r chi.Router) {
			r.Get("/", GetAuthors)
			r.Get("/{id}", GetAuthor)
		})
	})

	fmt.Println("Listening and Serving to ", Port)
	err := http.ListenAndServe(":"+Port, r)
	if err != nil {
		log.Fatal(err)
		return
	}

}
