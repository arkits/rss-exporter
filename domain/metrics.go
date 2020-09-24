package domain

import (
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

func MetricsHandler(w http.ResponseWriter, req *http.Request) {
	metrics.WritePrometheus(w, true)
}
