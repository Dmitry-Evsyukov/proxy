package scannerRepo

import (
	"database/sql"
	"encoding/json"
	"main/internal/models"
	"main/internal/pkg/scanner"
)

type Postgres struct {
	conn *sql.DB
}

func (p Postgres) GetAllRequests() ([]models.Request, error) {
	result := make([]models.Request, 0)
	query := `select id, method, scheme, url, headers, cookies, get_params, post_params from request`

	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		req := models.NewRequest()
		headers := make([]byte, 0)
		cookies := make([]byte, 0)
		getParams := make([]byte, 0)
		postParams := make([]byte, 0)

		err = rows.Scan(&req.Id, &req.Method, &req.Scheme, &req.Url, &headers, &cookies, &getParams, &postParams)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(headers, &req.Headers)
		json.Unmarshal(cookies, &req.Cookies)
		json.Unmarshal(getParams, &req.GetParams)
		json.Unmarshal(postParams, &req.PostParams)

		result = append(result, req)
	}
	return result, nil
}

func (p Postgres) GetRequest(id int) (models.Request, error) {
	req := models.NewRequest()
	headers := make([]byte, 0)
	cookies := make([]byte, 0)
	getParams := make([]byte, 0)
	postParams := make([]byte, 0)

	query := `select id, method, scheme, url, headers, cookies, get_params, post_params from request where id = $1`
	err := p.conn.QueryRow(query, id).Scan(&req.Id, &req.Method, &req.Scheme, &req.Url, &headers, &cookies, &getParams, &postParams)
	if err != nil {
		return models.Request{}, err
	}

	json.Unmarshal(headers, &req.Headers)
	json.Unmarshal(cookies, &req.Cookies)
	json.Unmarshal(getParams, &req.GetParams)
	json.Unmarshal(postParams, &req.PostParams)

	return req, nil
}

func (p Postgres) GetResponse(id int) (models.Response, error) {
	query := `select id, request_id, code, message, headers, body from response where request_id = $1`

	res := models.NewResponse()
	headers := make([]byte, 0)
	err := p.conn.QueryRow(query, id).Scan(&res.Id, &res.ReqId, &res.Code, &res.Message, &headers, &res.Body)
	if err != nil {
		return models.Response{}, err
	}

	json.Unmarshal(headers, &res.Headers)
	return res, nil
}

func New(c *sql.DB) scanner.Repository {
	return &Postgres{conn: c}
}
