package main

import (
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	//endpoint to check if server is ready to receive traffic
	w.WriteHeader(http.StatusOK)
	//writes a status code 200 (StatusOK) to ensure it works
	w.Write([]byte(http.StatusText(http.StatusOK)))
	//writes the data to the connection as part of an HTTP reply
}