package main

import "net/http"

func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isSafeMethod(r.Method) {
			setCSRFToken(w, r)
			next.ServeHTTP(w, r)
			return
		}
	})
}

func setCSRFToken(w http.ResponseWriter, r *http.Request) {

}

func isSafeMethod(method string) bool {
	return method == "GET" || method == "HEAD" || method == "OPTION"
}
