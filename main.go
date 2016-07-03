package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/v1/ping", ping)
	http.ListenAndServe(":9090", nil)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong!"))
}
