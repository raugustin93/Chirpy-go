package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/raugustin93/Chirpy-go/internal/db"
)

type apiConfig struct {
	fileserverHits int
	DB             *db.DB
	JwtSecret      []byte
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT Secret is not set")
	}

	const port = "8080"
	const filepathRoot = "."

	fileServerHandler := http.FileServer(http.Dir(filepathRoot))

	DB, err := db.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := apiConfig{
		fileserverHits: 0,
		DB:             DB,
		JwtSecret:      []byte(secretKey),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServerHandler)))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.HandlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGet)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", cfg.HandleLogin)
	mux.HandleFunc("PUT /api/users", cfg.HandlerUsersUpdate)
	mux.HandleFunc("POST /api/refresh", cfg.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.HandlerRevokeToken)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.HandlerChirpDelete)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server is listening on port 8080...")
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
