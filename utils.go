package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"html/template"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles("./templates/" + tmplName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func hashPassWord(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPassWordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Authorize(r *http.Request) error {
	authError := errors.New("Unauthorized Request")

	email := r.FormValue("email")
	user, ok := Users[email]
	if !ok {
		panic(errors.New("User does not exist"))
	}

	sToken, err := r.Cookie("session-token")
	if err != nil || sToken.Value == "" || sToken.Value != user.SessionToken {
		return authError
	}

	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken != user.CSRFToken || csrfToken == "" {
		return authError
	}
	return nil
}


