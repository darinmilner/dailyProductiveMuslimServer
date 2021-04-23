package main

import (
	"fmt"
	"log"
	"net/http"
	"server/everydaymuslimappserver/internal/helpers"

	"github.com/justinas/nosurf"
)

//WriteToConsole middleware--USELESS
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

//NoSurf adds CSRF to all POST requests
func NoSurf(next http.Handler) http.Handler {
	fmt.Println("NO SURF")
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad middleware loads and saves the session on each request
func SessionLoad(next http.Handler) http.Handler {
	fmt.Println("Session load")
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Auth Handler")
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Must be logged in!")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
