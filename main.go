package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
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
	serveMux := http.NewServeMux()
	const port = "8080"
	const filepathRoot = "."
	const apiPrefix = "/api/"
	const appPrefix = "/app/"
	const adminPrefix = "/admin/"

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}

	serveMux.Handle(appPrefix, apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc(createApiPath("GET", apiPrefix, "healthz"), hanlderHealth)
	serveMux.HandleFunc(createApiPath("GET", adminPrefix, "metrics"), apiCfg.handlerMetrics)
	serveMux.HandleFunc(createApiPath("POST ", adminPrefix, "reset"), apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	log.Fatal(server.ListenAndServe())
}
