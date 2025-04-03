package main

import (
	"crypto/rand"
	"encoding/base64"
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
		if isSafeMethod(r.Method) {
			setCSRFToken(w, r)
			next.ServeHTTP(w, r)
			return
		}

		if !validateCSRFToken(r) {
			http.Error(w, "CSRF token validation failed", http.StatusForbidden)
			return
		}

		setCSRFToken(w, r)
		next.ServeHTTP(w, r)
	})
}

func validateCSRFToken(r *http.Request) bool {
	cookieToken, err := r.Cookie(csrfCookieName)
	if err != nil {
		return false
	}

	requestToken := r.Header.Get(csrfHeaderName)
	if requestToken == "" {
		requestToken = r.FormValue(csrfFormFieldName)
	}
	return cookieToken.Value == requestToken
}

func setCSRFToken(w http.ResponseWriter, r *http.Request) {
	token, err := generateCSRFToken()
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		MaxAge:   3600, // 1 hour
		Path:     "/",
		Secure:   true,  // HTTPS only
		HttpOnly: false, // Allow JavaScript to read
		SameSite: http.SameSiteStrictMode,
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
	return base64.URLEncoding.EncodeToString(token), nil
}
