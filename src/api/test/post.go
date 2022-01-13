package test

import "net/http"

func Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post message received"))
}
