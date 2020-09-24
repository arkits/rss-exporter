package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arkits/rss-exporter/domain"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	version string
)

func init() {

	// Setup the Application wide config through Viper
	SetupConfig()

	go domain.BeginPollingFeeds()
}

func main() {
	port := viper.GetString("server.port")
	serviceName := viper.GetString("server.name")

	r := mux.NewRouter()

	r.HandleFunc("/", domain.VersionHandler).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s", serviceName), domain.VersionHandler).Methods(http.MethodGet)

	r.HandleFunc(fmt.Sprintf("/%s/feed", serviceName), domain.FeedHandler).Methods(http.MethodGet)

	// Expose Prometheus metrics
	r.HandleFunc(fmt.Sprintf("/%s/metrics", serviceName), domain.MetricsHandler).Methods(http.MethodGet)

	r.Use(domain.LoggingMiddleware)
	r.Use(domain.MetricsMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	log.Printf("Starting server on http://localhost:%v/%v", port, serviceName)
	http.ListenAndServe(":"+port, r)
}

// SetupConfig -  Setup the application config by reading the config file via Viper
func SetupConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file! - %s", err)
	}

	if version == "" {
		version = "undefined"
	}

	viper.Set("server.version", version)

}
