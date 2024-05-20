package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	fileServerHandler := http.FileServer(http.Dir(filepathRoot))
	cfg := apiConfig{fileserverHits: 0}

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServerHandler)))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.HandlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server is listening on port 8080...")
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
