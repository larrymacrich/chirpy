package server

import (
	"log"
	"net/http"
	"sync/atomic"

	api "github.com/larrymacrich/chirpy/internal/api"
	fileserver "github.com/larrymacrich/chirpy/internal/frontend"
	metrics "github.com/larrymacrich/chirpy/internal/metrics"
)

func StartServer(port string) {
	metricsCfg := metrics.MetricsConfig{
		FileserverRequests: atomic.Int32{},
	}

	fsRoot := "."
	fsPrefix := "/app"

	mux := http.NewServeMux()

	// fileserver endpoint (serves static files, tracked via metrics)
	mux.Handle("/app/", fileserver.HandlerFileserver(&metricsCfg, fsPrefix, fsRoot))

	// health check
	mux.HandleFunc("GET /api/healthz", api.HandlerReadiness)

	// metrics reporting/reset
	mux.HandleFunc("GET /api/metrics", metricsCfg.HandlerMetricsRequest)
	mux.HandleFunc("POST /api/reset", metricsCfg.HandlerMetricsReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", fsRoot, port)
	log.Fatal(srv.ListenAndServe())
}
