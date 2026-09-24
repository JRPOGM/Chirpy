package main

import "net/http"

func (cfg *apiConfig) handlerResets(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	cfg.fileserverHits.Store(0)
	// atomic.Int23 .Store(val int32) method stores a value to call and enforce
	err := cfg.db.Reset(r.Context())
	// cfg.db.Reset calls new query function stored in reset.sql.go
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}