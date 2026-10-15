package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	filepathRoot = "."
	port         = "8080"
)

func main() {

	startServer()
}

func startServer() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	apiCfg.fileserverHits.Store(0)

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", apiCfg.handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerNumberRequest)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return cfg.handlerLogHits(next)
}

func (cfg *apiConfig) handlerLogHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerNumberRequest(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
