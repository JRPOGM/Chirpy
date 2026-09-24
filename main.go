package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"

	"github.com/JRPOGM/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	db				*database.Queries
	platform		string
} //holds stateful, in-memory data to keep track of
//atomic.Int32 allows to safely incriment & read int value across multiple goroutines

func main() {
	godotenv.Load()
	//loads the .env file into environment variables
	dbURL := os.Getenv("DB_URL")
	//os.Getenv(key string) gets database url from the environment
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	//sql.Open(driverName, dataSourceName string) opens a connection to your database
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)
	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform,
	}
	multiplex := http.NewServeMux()
	//routes requests and creates a simple http server for the program
	multiplex.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	//.Handle(string, Handler) adds a handler for the path of the root directory
	//http.StripPrefix(string, Handler) strips the prefix from the request path before passing it to the fileserver handler
	//http.FileServer(http.Dir(root file path)) converts a filepath to a directory to use as the Handler
	multiplex.HandleFunc("GET /api/healthz", handlerReadiness)
	// .HandleFunc([host]/[path]) setup to specify methods for functions
	multiplex.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	multiplex.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	// /api path for decoupled logic
	// non GET functions should return a 405 status response (Method Not Allowed)
	multiplex.HandleFunc("POST /admin/reset", apiCfg.handlerResets)
	multiplex.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// /admin path are for user control data
	//apiCfg.func prefix for all functions with (cfg *apiConfig)
	serve := &http.Server{
		Addr:	":" + port,
		Handler: multiplex,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serve.ListenAndServe())
	//starts the server
}