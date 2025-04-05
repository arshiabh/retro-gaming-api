package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

const (
	csrfTokenLength   = 32 // 32 bytes = 256 bits
	csrfCookieName    = "csrf_token"
	csrfHeaderName    = "X-CSRF-Token"
	csrfFormFieldName = "csrf_token"
)

func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Cookies())
		if isSafeMethod(r.Method) {
			setCSRFToken(w)
			next.ServeHTTP(w, r)
			return
		}

		if !validateCSRFToken(r) {
			http.Error(w, "CSRF token validation failed", http.StatusForbidden)
			return
		}

		setCSRFToken(w)
		next.ServeHTTP(w, r)
	})
}

func validateCSRFToken(r *http.Request) bool {
	cookieToken, err := r.Cookie("csrf_token")
	if err != nil {
		log.Println(err)
		return false
	}

	requestToken := r.Header.Get(csrfHeaderName)
	if requestToken == "" {
		requestToken = r.FormValue(csrfFormFieldName)
	}
	
	return cookieToken.Value == requestToken
}

func setCSRFToken(w http.ResponseWriter) {
	token, err := generateCSRFToken()
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		MaxAge:   3600, // 1 hour
		Path:     "/",
		Secure:   false, // HTTPS only
		HttpOnly: false, // Allow JavaScript to read
		SameSite: http.SameSiteLaxMode,
	})
}

func isSafeMethod(method string) bool {
	return method == "GET" || method == "HEAD" || method == "OPTION"
}

func generateCSRFToken() (string, error) {
	token := make([]byte, csrfTokenLength)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	fmt.Println(base64.RawURLEncoding.EncodeToString(token))
	return base64.RawURLEncoding.EncodeToString(token), nil
}
