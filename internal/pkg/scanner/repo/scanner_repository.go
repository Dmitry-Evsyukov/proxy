package scannerRepo

import (
	"database/sql"
	"main/internal/models"
	"main/internal/pkg/scanner"
)

type Postgres struct {
	conn *sql.DB
}

func (p Postgres) GetAllRequests() ([]models.Request, error) {
	result := make([]models.Request, 0)
	query := `select id, data from request`

	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var req models.Request
		req.Req = make([]byte, 0)

		err = rows.Scan(&req.Id, &req.Req)
		if err != nil {
			return nil, err
		}

		result = append(result, req)
	}
	return result, nil
}

func (p Postgres) GetRequest(id int) (models.Request, error) {
	var result models.Request
	query := `select id, data from request where id = $1`

	result.Req = make([]byte, 0)
	err := p.conn.QueryRow(query, id).Scan(&result.Id, &result.Req)
	if err != nil {
		return models.Request{}, err
	}

	return result, nil
}

func (p Postgres) GetResponse(id int) (models.Response, error) {
	var result models.Response
	query := `select id, data from response where id = $1`

	result.Res = make([]byte, 0)
	err := p.conn.QueryRow(query, id).Scan(&result.Id, &result.Res)
	if err != nil {
		return models.Response{}, err
	}

	return result, nil
}

func New(c *sql.DB) scanner.Repository {
	return &Postgres{conn: c}
}
