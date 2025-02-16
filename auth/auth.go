package auth

import (
	"encoding/json"
	"fmt"
	"github.com/RejwankabirHamim/api-book-server/data"
	"net/http"
	"time"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	var cred data.Credential
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}
	expectedPassword, ok := data.Users[cred.Username]
	fmt.Println(expectedPassword, ok)

	if !ok || expectedPassword != cred.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(time.Hour * 2)
	_, tokenString, err := data.TokenAuth.Encode(map[string]interface{}{
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

func LogOut(w http.ResponseWriter, r *http.Request) {
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
