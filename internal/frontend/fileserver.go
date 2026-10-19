package fileserver

import (
	"net/http"

	metrics "github.com/larrymacrich/chirpy/internal/metrics"
)

func HandlerFileserver(cfg *metrics.MetricsConfig, fsroot string, prefix string) http.Handler {
	return cfg.MiddlewareMetricsInc(http.StripPrefix(prefix, http.FileServer(http.Dir(fsroot))))
}
