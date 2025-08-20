package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/federicoReghini/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	template := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

	fmt.Fprintf(w, template, c.fileServerHits.Load())
}

func createApiPath(method, prefix, path string) string {
	return fmt.Sprintf("%s %s%s", method, prefix, path)
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	const port = "8080"
	const filepathRoot = "."
	const apiPrefix = "/api/"
	const appPrefix = "/app/"
	const adminPrefix = "/admin/"

	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
	}

	serveMux.Handle(appPrefix, apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc(createApiPath("GET", apiPrefix, "healthz"), handlerHealth)
	serveMux.HandleFunc(createApiPath("GET", adminPrefix, "metrics"), apiCfg.handlerMetrics)
	serveMux.HandleFunc(createApiPath("POST ", adminPrefix, "reset"), apiCfg.handlerReset)
	serveMux.HandleFunc(createApiPath("POST ", apiPrefix, "chirps"), apiCfg.handlerCreateChirp)
	serveMux.HandleFunc(createApiPath("POST ", apiPrefix, "users"), apiCfg.handlerCreateUser)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	log.Fatal(server.ListenAndServe())
}
