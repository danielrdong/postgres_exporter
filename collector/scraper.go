package collector

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
)

type Scraper interface {
	// Name of the Scraper. Should be unique.
	Name() string
	// Scrape collects data from database connection and sends it over channel as prometheus metric.
	Scrape(db *sql.DB, ch chan<- prometheus.Metric) error
}