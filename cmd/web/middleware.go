package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf" // New import
)

type ContextKey string

const LIC_INFO ContextKey = "LIC_INFO"

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {

	defaultFailureHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println(" :::::::::::: CSRF FAILED ::::::::::::::::", nosurf.Reason(r))
		http.Error(w, http.StatusText(400), 400)
	})

	csrfHandler := nosurf.New(next)
	// csrfHandler.SetBaseCookie(http.Cookie{
	// 	HttpOnly: true,
	// 	//Path:     "/",
	// 	//Secure: true,
	// })
	csrfHandler.SetFailureHandler(defaultFailureHandler)
	return csrfHandler
}

const (
	xForwardedProtoHeader = "x-forwarded-proto"
)

func (app *application) RedirectToHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//host, _, _ := net.SplitHostPort(r.Host)
		u := r.URL
		log.Println("starte", u.String(), "::", r.URL.Scheme, r.TLS, r.Host, r.RequestURI, "::", r.Header.Get(xForwardedProtoHeader))
		if r.Header.Get(xForwardedProtoHeader) != "https" {

			//log.Println(":::::::: REDIRECTING :::::::::")
			sslUrl := "https://" + r.Host + r.RequestURI
			http.Redirect(w, r, sslUrl, http.StatusMovedPermanently)
			return
		}

		//log.Println(":::::::: NOT REDIRECTING :::::::::")

		next.ServeHTTP(w, r)
	})
}
