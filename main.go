package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ratludu/http-servers-go/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwt            string
	polka          string
}

func main() {

	godotenv.Load()

	plt := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	JWT := os.Getenv("JWT_KEY")
	Polka := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       plt,
		jwt:            JWT,
		polka:          Polka,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/heathz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpDelete)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersPut)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgrade)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
