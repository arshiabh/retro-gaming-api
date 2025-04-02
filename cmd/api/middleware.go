package main

import "net/http"

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

}

func isSafeMethod(method string) bool {
	return method == "GET" || method == "HEAD" || method == "OPTION"
}
