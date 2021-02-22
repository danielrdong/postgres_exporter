package scrapers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/danielrdong/postgres_exporter/collector"
	"github.com/danielrdong/postgres_exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	maxConnSql = `show max_connections;`
	curConnSql = `select count(*) from pg_stat_activity;`
	actConnSql = `select count(1) from pg_stat_activity where not pid=pg_backend_pid();`
)

type PgConnScraper struct{}

func (PgConnScraper) Name() string {
	return "postgres_connections_scraper"
}

func (PgConnScraper) Scrape(db *sql.DB, ch chan<- prometheus.Metric) error {
	maxConn, errM := queryMaxConn(db)
	curConn, errC := queryCurConn(db)
	ActConn, errA := queryActConn(db)

	ch <- prometheus.MustNewConstMetric(maxConnDesc, prometheus.GaugeValue, maxConn)
	ch <- prometheus.MustNewConstMetric(curConnDesc, prometheus.GaugeValue, curConn)
	ch <- prometheus.MustNewConstMetric(actConnDesc, prometheus.GaugeValue, ActConn)

	return utils.CombineErr(errM, errC, errA)
}

func queryMaxConn(db *sql.DB) (maxConn float64, err error) {
	rows, err := db.Query(maxConnSql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&maxConn)
		return
	}
	err = errors.New(fmt.Sprintf("%s not found.", maxConnSql))
	return
}
func queryCurConn(db *sql.DB) (curConn float64, err error) {
	rows, err := db.Query(curConnSql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&curConn)
		return
	}
	err = errors.New(fmt.Sprintf("%s not found.", curConnSql))
	return
}
func queryActConn(db *sql.DB) (actConn float64, err error) {
	rows, err := db.Query(actConnSql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&actConn)
		return
	}
	err = errors.New(fmt.Sprintf("%s not found.", actConnSql))
	return

}

var (
	maxConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, collector.SystemExporter, "max_connections"),
		"Get the metrics of the max connections of Postgres.",
		nil,
		nil,
	)

	curConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, collector.SystemExporter, "current_connections"),
		"Get the metrics of the current connections of Postgres.",
		nil,
		nil,
	)

	actConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, collector.SystemExporter, "active_connections"),
		"Get the metrics of the active connections of Postgres.",
		nil,
		nil,
	)
)
