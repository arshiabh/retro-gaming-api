package main

import "net/http"

func (app *application) HandleGetTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello word"))
}
