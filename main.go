package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var Users = make(map[string]*User)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fs))

	//static assets (CSS, JS, etc.)
	staticFS := http.FileServer(http.Dir("templates/styles"))
	mux.Handle("/styles/", http.StripPrefix("/styles/", staticFS))

	// For images in templates folder
	imagesFS := http.FileServer(http.Dir("templates/images"))
	mux.Handle("/images/", http.StripPrefix("/images/", imagesFS))

	//pages
	mux.HandleFunc("/", HomePage)
	mux.HandleFunc("/permits", PermitsPage)
	mux.HandleFunc("/signup", SignUpPage)
	mux.HandleFunc("/signinp", SignInPage)

	mux.HandleFunc("/applications/zone-license", ZoneLicense)
	mux.HandleFunc("/applications/trading-license", TradingLicense)
	mux.HandleFunc("/applications/fire-clearence", FireClearence)

	//requests
	mux.HandleFunc("POST /api/signup", HandleSignUp)
	mux.HandleFunc("POST /api/signin", HandleSignIn)

	log.Printf("Server listening on port :8000")
	http.ListenAndServe(":8000", mux)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func PermitsPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Permits.html", nil)
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Signup.html", nil)
}

func SignInPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login.html", nil)
}

func ZoneLicense(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./applications/ZoneLicence.html", nil)
}

func TradingLicense(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./applications/TradingLicense.html", nil)
}

func FireClearence(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./applications/FireClearence.html", nil)
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword, err := hashPassWord(password)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusUnauthorized)
	}

	if _, ok := Users[email]; ok {
		http.Error(w, "User already exists", http.StatusUnauthorized)
	}

	user := &User{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	Users[email] = user
	fmt.Fprintf(w, "User %s created successfully", user.Name)
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, ok := Users[email]
	if !ok {
		http.Error(w, "Invalid Email Address", http.StatusUnauthorized)
		return
	}

	if !checkPassWordHash(user.HashedPassword, password) {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)
	CSRFToken := generateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf-token",
		Value:    CSRFToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	user.CSRFToken = CSRFToken
	user.SessionToken = sessionToken
	fmt.Fprintf(w, "Welcome %s", user.Name)

	//redirect user to Licenses list?
}
