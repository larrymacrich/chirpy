package metrics

import (
	"fmt"
	"net/http"
)

func (cfg *MetricsConfig) HandlerMetricsRequest(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(createMetricsHTMLContent(cfg)))
}

func createMetricsHTMLContent(cfg *MetricsConfig) string {
	return fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileserverRequests.Load())
}

func (cfg *MetricsConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return cfg.handlerMetricsInc(next)
}

func (cfg *MetricsConfig) handlerMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverRequests.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *MetricsConfig) HandlerMetricsReset(w http.ResponseWriter, req *http.Request) {
	cfg.FileserverRequests.Store(0)
	w.WriteHeader(http.StatusOK)
}
