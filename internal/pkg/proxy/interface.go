package proxy

import "net/http"

type Repository interface {
	SaveResponse(response *http.Response) error
	SaveRequest(request *http.Request) error
}
