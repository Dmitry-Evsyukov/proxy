package proxyRepo

import (
	"database/sql"
	"main/internal/pkg/proxy"
	"main/internal/pkg/utils"
	"net/http"
)

type Postgres struct {
	conn *sql.DB
}

func (p *Postgres) SaveRequest(req *http.Request) error {
	data, err := utils.RequestToJson(req)
	if err != nil {
		return err
	}

	query := `insert into request (data) values ($1)`
	_, err = p.conn.Exec(query, data)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) SaveResponse(resp *http.Response) error {
	data, err := utils.ResponseToJson(resp)
	if err != nil {
		return err
	}

	query := `insert into response (data) values ($1)`
	_, err = p.conn.Exec(query, data)
	if err != nil {
		return err
	}

	return nil
}

func New(db *sql.DB) proxy.Repository {
	return &Postgres{conn: db}
}
