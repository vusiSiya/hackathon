package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"
)

var Users = make(map[string]*User)
var Applications = make(map[string]*BusinessApplication)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fs))

	//static assets
	staticFS := http.FileServer(http.Dir("templates/styles"))
	mux.Handle("/styles/", http.StripPrefix("/styles/", staticFS))

	//images
	imagesFS := http.FileServer(http.Dir("templates/images"))
	mux.Handle("/images/", http.StripPrefix("/images/", imagesFS))

	//pages
	mux.HandleFunc("/", HomePage)
	mux.HandleFunc("/index", IndexPage)
	mux.HandleFunc("/permits", PermitsPage)
	mux.HandleFunc("/signup", SignUpPage)
	mux.HandleFunc("/login", SignInPage)
	mux.HandleFunc("POST /signout", SignOut)

	mux.HandleFunc("/applications/zone-license", ZoneLicense)
	mux.HandleFunc("/applications/trading-license", TradingLicense)
	mux.HandleFunc("/applications/fire-clearence", FireClearence)

	//requests
	mux.HandleFunc("POST /api/signup", HandleSignUp)
	mux.HandleFunc("POST /api/login", HandleSignIn)
	mux.HandleFunc("POST /api/apply-zone-license", HandleApplyZoneLicense)

	log.Printf("Server listening on port :8000")
	http.ListenAndServe(":8000", mux)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "WelcomePage.html", nil)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func PermitsPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Permits.html", nil)
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "signup.html", nil)
}

func SignInPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
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
	log.Printf("User %s created successfully", user.Name)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	log.Printf("User %s signed in successfully", user.Name)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleApplyZoneLicense(w http.ResponseWriter, r *http.Request) {
	log.Printf("Applying for Zone License")

	err := Authorize(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Unauthorized attempt to apply for Zone License: %v", err)
		return
	}

	application := &BusinessApplication{
		ApplicationID:   rand.Text(),
		UserEmail:       r.FormValue("email"),
		ApplicationType: r.FormValue("application_type"),
		BusinessName:    r.FormValue("business_name"),
		BusinessType:    r.FormValue("business_type"),
		Description:     r.FormValue("description"),
		Address:         r.FormValue("address"),
		DateCreated:     time.Now(),
	}
	Applications[application.ApplicationID] = application
	fmt.Fprintf(w, "Application for %s submitted successfully", application.ApplicationType)
	log.Printf("Application for %s submitted successfully", application.ApplicationType)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	// Just use the session token to identify the user
	sessionToken, err := r.Cookie("session-token")
	if err != nil {
		// Already logged out or no session - that's fine
		fmt.Fprintln(w, "Signed out")
		return
	}

	// Find user by their session token, not by form input
	user, ok := findUserBySessionToken(sessionToken.Value)
	if ok {
		user.SessionToken = ""
		user.CSRFToken = ""
	}

	// Clear cookies regardless
	clearCookie(w, "session-token", true)
	clearCookie(w, "csrf-token", false)
}

func clearCookie(w http.ResponseWriter, name string, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: httpOnly,
		Secure:   true, // Always set in production
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}
