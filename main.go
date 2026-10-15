package main

import (
	"log"
	"net/http"
)

const (
	filepathRoot = "."
	port         = "8080"
)

func main() {
	startServer()
}

func startServer() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)
	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
