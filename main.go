package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	multiplex := http.NewServeMux()
	multiplex.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	multiplex.HandleFunc("GET /healthz", handlerReadiness)
	multiplex.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	multiplex.HandleFunc("POST /reset", apiCfg.handlerResets)
	serve := &http.Server{
		Addr:	":" + port,
		Handler: multiplex,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serve.ListenAndServe())
}

