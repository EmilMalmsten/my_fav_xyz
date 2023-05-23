package main

import "net/http"

func (cfg apiConfig) handlerCreateToplist(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test123"))
}
