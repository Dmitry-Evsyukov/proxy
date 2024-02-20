package proxyRepo

import (
	"database/sql"
	"main/internal/pkg/proxy"
	"net/http"
)

type Postgres struct {
	conn *sql.DB
}

func (p *Postgres) SaveResponse(resp *http.Response) error {
	return nil
}

func (p *Postgres) SaveRequest(req *http.Request) error {
	return nil
}

func New(db *sql.DB) proxy.Repository {
	return &Postgres{conn: db}
}
