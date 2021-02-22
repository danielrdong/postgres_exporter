package collector

import "github.com/prometheus/client_golang/prometheus"

const (
	Namespace      = "postgres"
	SystemExporter = "exporter"
)

type Metrics struct {
	totalScraped   prometheus.Counter
	totalError     prometheus.Counter
	scrapeDuration prometheus.Gauge
	postgresUp     prometheus.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		totalScraped: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: SystemExporter,
				Name:      "total_scraped",
				Help:      "Total scraped",
			},
		),
		totalError: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: SystemExporter,
				Name:      "total_error",
				Help:      "Total error scraping",
			},
		),
		scrapeDuration: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: SystemExporter,
				Name:      "scrape_duration_second",
				Help:      "Elapsed of each scrape",
			},
		),
		postgresUp: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "up",
				Help:      "Whether postgres is reachable",
			},
		),
	}
}
