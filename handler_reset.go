package main

import "net/http"

func (cfg *apiConfig) handlerResets(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	// atomic.Int23 .Store(val int32) method stores a value to call and enforce
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}