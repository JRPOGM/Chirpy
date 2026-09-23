package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
} //holds stateful, in-memory data to keep track of
//atomic.Int32 allows to safely incriment & read int value across multiple goroutines

func main() {
	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	multiplex := http.NewServeMux()
	//routes requests and creates a simple http server for the program
	multiplex.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	//.Handle(string, Handler) adds a handler for the path of the root directory
	//http.StripPrefix(string, Handler) strips the prefix from the request path before passing it to the fileserver handler
	//http.FileServer(http.Dir(root file path)) converts a filepath to a directory to use as the Handler
	multiplex.HandleFunc("GET /api/healthz", handlerReadiness)
	multiplex.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	multiplex.HandleFunc("POST /admin/reset", apiCfg.handlerResets)
	multiplex.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serve := &http.Server{
		Addr:	":" + port,
		Handler: multiplex,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serve.ListenAndServe())
	//starts the server
}