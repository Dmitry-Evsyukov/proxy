package init_db

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

func NewConn(url string) *sql.DB {
	result, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalln("error conn db", err)
	}

	return result
}
