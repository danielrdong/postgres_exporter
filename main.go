package main

import (
	"github.com/danielrdong/postgres_exporter/collector"
	demo "github.com/danielrdong/postgres_exporter/scrapers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	listenAddress, metricsPath = "0.0.0.0:9297", "/metrics"
	disableDefaultMetrics      = true

	gathers prometheus.Gatherers
)

var scrapers = map[collector.Scraper]bool{
	demo.PgConnScraper{}: true,
}

func main() {
	metricsHandleFunc := newHandler(disableDefaultMetrics, scrapers)

	mux := http.NewServeMux()

	mux.HandleFunc(metricsPath, metricsHandleFunc)
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}

func newHandler(disableDefaultMetrics bool, scrapers map[collector.Scraper]bool) http.HandlerFunc {
	registry := prometheus.NewRegistry()

	enabledScrapers := make([]collector.Scraper, 0)

	for scraper, enable := range scrapers {
		if enable {
			enabledScrapers = append(enabledScrapers, scraper)
		}
	}

	postgresCollector := collector.NewCollector(enabledScrapers)
	registry.MustRegister(postgresCollector)

	if disableDefaultMetrics {
		gathers = prometheus.Gatherers{registry}
	} else {
		gathers = prometheus.Gatherers{registry, prometheus.DefaultGatherer}
	}

	handler := promhttp.HandlerFor(gathers, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})

	return handler.ServeHTTP
}
