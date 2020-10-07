package domain

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

// LoggingMiddleware logs incoming HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, requestToIgnoredEndpoint := Find(GetIgnoredEndpoints(), r.RequestURI)

		if !requestToIgnoredEndpoint {
			log.Printf("%s - %s | %s", r.Method, r.RequestURI, r.RemoteAddr)
		}

		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware logs metrics related to the HTTP requests and responses
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func(begin time.Time) {

			requestDuration := fmt.Sprintf(`http_requests_duration_seconds{path="%v", method="%v"}`,
				r.RequestURI, r.Method,
			)
			metrics.GetOrCreateSummary(requestDuration).UpdateDuration(begin)

		}(time.Now())

		next.ServeHTTP(w, r)
	})
}

func GetIgnoredEndpoints() []string {
	return []string{"/rss/metrics"}
}
