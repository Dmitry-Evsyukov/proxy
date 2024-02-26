package proxy

import "net/http"

type Repository interface {
	SaveResponse(response *http.Response, reqId int) error
	SaveRequest(request *http.Request) (int, error)
}
