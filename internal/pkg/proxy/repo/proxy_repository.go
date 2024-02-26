package proxyRepo

import (
	"database/sql"
	"encoding/json"
	"main/internal/pkg/proxy"
	"main/internal/pkg/utils"
	"net/http"
)

type Postgres struct {
	conn *sql.DB
}

func (p *Postgres) SaveRequest(req *http.Request) (int, error) {
	reqStruct, err := utils.RequestToStruct(req)
	if err != nil {
		return 0, err
	}

	headers, err := json.Marshal(reqStruct.Headers)
	if err != nil {
		return 0, err
	}

	cookies, err := json.Marshal(reqStruct.Cookies)
	if err != nil {
		return 0, err
	}

	getParams, err := json.Marshal(reqStruct.GetParams)
	if err != nil {
		return 0, err
	}

	postParams, err := json.Marshal(reqStruct.PostParams)
	if err != nil {
		return 0, err
	}

	var id int
	query := `INSERT INTO request (method, scheme, url, headers, cookies, get_params, post_params) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = p.conn.QueryRow(query, reqStruct.Method, reqStruct.Scheme, reqStruct.Url, headers, cookies, getParams, postParams).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

//request_id int references request(id),
//code int,
//message text,
//headers json,
//body text

func (p *Postgres) SaveResponse(resp *http.Response, reqId int) error {
	respStruct, err := utils.ResponseToStruct(resp, reqId)
	if err != nil {
		return err
	}

	headers, err := json.Marshal(respStruct.Headers)
	if err != nil {
		return err
	}

	query := `insert into response (request_id, code, message, headers, body) values ($1, $2, $3, $4, $5)`
	_, err = p.conn.Exec(query, reqId, respStruct.Code, respStruct.Message, headers, respStruct.Body)
	if err != nil {
		return err
	}

	return nil
}

func New(db *sql.DB) proxy.Repository {
	return &Postgres{conn: db}
}
