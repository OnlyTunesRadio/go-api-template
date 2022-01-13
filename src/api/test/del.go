package test

import "net/http"

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Item deleted successfully."))
}
