package scannerRepo

import (
	"database/sql"
	"main/internal/pkg/scanner"
)

type Postgres struct {
	conn *sql.DB
}

func (p Postgres) GetAllRequests() ([][]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) GetRequest(id int) ([]byte, error) {
	query := `select data from request where id = $1`

	data := make([]byte, 0)
	err := p.conn.QueryRow(query, id).Scan(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p Postgres) GetResponse(id int) ([]byte, error) {
	query := `select data from response where id = $1`

	data := make([]byte, 0)
	err := p.conn.QueryRow(query, id).Scan(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func New(c *sql.DB) scanner.Repository {
	return &Postgres{conn: c}
}
