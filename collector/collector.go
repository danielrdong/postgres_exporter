package collector

import (
	"database/sql"
	"github.com/danielrdong/postgres_exporter/utils"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"sync"
	"time"
)

const checkSql = `select 'OK';`

type PostgresCollector struct {
	mu       sync.Mutex
	db       *sql.DB
	metrics  *Metrics
	scrapers []Scraper
}

func NewCollector(enabledScrapers []Scraper) *PostgresCollector {
	return &PostgresCollector{
		metrics:  NewMetrics(),
		scrapers: enabledScrapers,
	}
}

func (pg *PostgresCollector) Collect(ch chan<- prometheus.Metric) {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	pg.scrape(ch)

	ch <- pg.metrics.totalScraped
	ch <- pg.metrics.totalError
	ch <- pg.metrics.scrapeDuration
	ch <- pg.metrics.postgresUp
}

func (pg *PostgresCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pg.metrics.postgresUp.Desc()
	ch <- pg.metrics.scrapeDuration.Desc()
	ch <- pg.metrics.totalError.Desc()
	ch <- pg.metrics.totalScraped.Desc()
}

func (pg *PostgresCollector) scrape(ch chan<- prometheus.Metric) {
	start := time.Now()

	//watch := stopwatch.New("scrape")

	pg.metrics.totalScraped.Inc()
	//watch.MustStart("check connections")
	err := pg.checkDBConn()
	//watch.MustStop()

	if err != nil {
		pg.metrics.totalError.Inc()
		pg.metrics.scrapeDuration.Set(time.Since(start).Seconds())
		pg.metrics.postgresUp.Set(0)

		log.Error("check database connection failed. error: %v", err.Error())
		return
	}

	pg.metrics.postgresUp.Set(1)

	for _, scraper := range pg.scrapers {
		//watch.MustStart("scraping: " + scraper.Name())
		err := scraper.Scrape(pg.db, ch)
		//watch.MustStop()
		if err != nil {
			log.Error("get metrics for scraper: %s failed, error: %v", scraper.Name(), err.Error())
		}
	}

	pg.metrics.scrapeDuration.Set(time.Since(start).Seconds())
	//log.Info(fmt.Sprintf("prometheus scraped postgres exporter successfully at %v, datail elapsed: %s", time.Now(), watch.PrettyPrint()))
}

func (pg *PostgresCollector) checkDBConn() error {
	if pg.db == nil {
		return pg.getDBConn()
	}

	if err := checkPGConnections(pg.db); err == nil {
		return nil
	} else {
		_ = pg.db.Close()
		pg.db = nil
		return pg.getDBConn()
	}
}

func (pg *PostgresCollector) getDBConn() error {

	db, err := sql.Open("postgres", utils.GetDbInfo())
	if err != nil {
		return err
	}

	if err = checkPGConnections(db); err != nil {
		_ = db.Close()
		return err
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	pg.db = db
	return nil
}

func checkPGConnections(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	rows, err := db.Query(checkSql)
	if err != nil {
		return err
	}

	defer func() {
		_ = rows.Close()
	}()
	return nil
}
