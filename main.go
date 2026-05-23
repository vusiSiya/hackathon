package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var Users = make(map[string]User)

func main() {

	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/zone-license", Zoning)
	mux.HandleFunc("/signup", Signup)

	mux.HandleFunc("/create-account", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hashedPassword, err := hashPassWord(password)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusBadRequest)
		}

		if _, ok := Users[email]; ok {
			http.Error(w, "User already exists", http.StatusBadRequest)
		}

		user := User{
			Name:           name,
			Email:          email,
			HashedPassword: hashedPassword,
		}
		Users[email] = user
		fmt.Fprintf(w, "User %s created successfully", user.Name)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// email := r.FormValue("email")
		// password := r.FormValue("password")

		// fmt.Printf(w, "User %s created successfully", user.Name)
	})

	log.Printf("Server listening on port :8000")

	http.ListenAndServe(":8000", mux)
}

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

func Zoning(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "application.html", nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "signup.html", nil)
}
