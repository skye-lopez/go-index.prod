// This is data service for the go-index API and static dist. The raw data is based on the https://index.golang.org/ api.
package idx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/skye-lopez/go-index.prod/pg"
)

// the package json structure from each line read in on the https://index.golang.org/index endpoint
type IdxEntry struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

// @FetchIdx - Given a properly configured postgresql instance it queries the https://index.golang.org api
// using the since param to ensure all packages are found and stored with their version information as well.
func FetchIdx() {
	db, err := pg.NewPG()
	if err != nil {
		log.Fatalf("Error establishing connection to db during FetchIdx();\n%s", err)
	}
	defer db.Conn.Close()

	startTime := time.Now()
	endTime, err := time.Parse(time.RFC3339Nano, getLastFetchTime(db.Conn))

	urls := generateUrls(startTime, endTime)
	var wg sync.WaitGroup
	sem := make(chan int, getMaxWorkers())
	entries := make(chan *IdxEntry, len(urls)*2000)

	for i, url := range urls {
		fmt.Printf("\r Getting package entries from urls => [ %d / %d ]", i, len(urls))

		sem <- 1
		wg.Add(1)
		// TODO: eventually theres should be some meaningful error handling here.
		// ideally only if there is a lot of non 200 status'
		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				<-sem
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				<-sem
				return
			}

			lines := strings.Split(string(body), "\n")
			for _, e := range lines {
				ie := &IdxEntry{}
				json.Unmarshal([]byte(e), ie)
				if len(ie.Path) < 5 {
					<-sem
					return
				}
				entries <- ie
			}
			<-sem
		}()
	}

	wg.Wait()
	close(entries)

	for e := range entries {
		fmt.Printf("\r Storing package entries to db =>  [ %d ] left", len(entries))
		sem <- 1
		wg.Add(1)

		// TODO: Db err handling. Ideally we could turn this into a tx and store failed ones to some kind of log to be processed later.
		go func() {
			defer wg.Done()

			// Upsert package into packages table
			db.Conn.Exec("INSERT INTO packages (url) VALUES ($1) ON CONFLICT DO NOTHING", e.Path)

			// Upsert package version for that package into package_versions
			r, _ := db.QueryString("SELECT EXISTS(SELECT version FROM package_versions WHERE owner = $1 and version = $2)",
				e.Path,
				e.Version,
			)
			exists := r[0].([]interface{})[0].(bool)

			if !exists {
				db.Conn.Exec("INSERT INTO package_versions (owner, version, time) VALUES ($1, $2, $3)",
					e.Path,
					e.Version,
					e.Timestamp,
				)
			}
			<-sem
		}()
	}

	wg.Wait()
	close(sem)

	newLastFetchTime := time.Now().Format(time.RFC3339Nano)
	db.Conn.Exec("UPDATE internal_log SET value = $1 WHERE id = $2",
		newLastFetchTime,
		"last_fetch_time",
	)
}

// gets the last time we fetched the index from the db. in RFC3339Nano format.
func getLastFetchTime(conn *sql.DB) string {
	// TODO: Remove this once all testing is done.
	devEnv := os.Getenv("DEV_ENV")
	if devEnv == "dev" {
		step := time.Duration(200) * time.Hour
		return time.Now().Add(-step).Format(time.RFC3339Nano)
	}

	var lastFetchTime string
	err := conn.QueryRow("SELECT value FROM internal_log WHERE id = $1", "last_fetch_time").Scan(&lastFetchTime)
	if err != nil {
		log.Fatalf("Error getting the last_write_time from db\n%s", err)
	}

	return lastFetchTime
}

// generates a list of all urls to fetch between the startTime and endTime(last_fetch_time from the db)
func generateUrls(startTime time.Time, endTime time.Time) []string {
	baseUrl := "https://index.golang.org/index?since="
	urls := []string{}
	step := time.Duration(12) * time.Hour
	for startTime.Unix() > endTime.Unix() {
		urls = append(urls, baseUrl+startTime.Format(time.RFC3339Nano))
		startTime = startTime.Add(-step)
	}

	return urls
}

// gets the maxWorkers for the given machine or returns a default of 10
func getMaxWorkers() int {
	maxWorkers := os.Getenv("MAX_WORKERS")
	workers, err := strconv.Atoi(maxWorkers)
	if err != nil {
		return 10
	}
	return workers
}
