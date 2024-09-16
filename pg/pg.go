package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	gq "github.com/skye-lopez/go-query"
)

func NewPG() (*gq.GoQuery, error) {
	connString := ""
	devEnv := os.Getenv("DEV_ENV")
	if devEnv == "dev" {
		connString = fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("PG_USER_DEV"),
			os.Getenv("PG_PWD_DEV"),
			os.Getenv("PG_DBNAME"),
			os.Getenv("PG_PORT"),
		)
	} else if devEnv == "prod" {
	} else {
		return &gq.GoQuery{}, errors.New("NewPG(); no devEnv provided, cannot create db connString.")
	}

	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return &gq.GoQuery{}, err
	}

	db := gq.NewGoQuery(conn)

	_, err = db.Conn.Exec("SELECT 1 as test")
	if err != nil {
		return &db, err
	}

	return &db, nil
}
