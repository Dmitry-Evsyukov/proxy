package models

type Request struct {
	Id         int
	Method     string
	Scheme     string
	Url        string
	GetParams  map[string]any
	Headers    map[string]any
	Cookies    map[string]any
	PostParams map[string]any
}

type Response struct {
	Id      int
	ReqId   int
	Code    int
	Message string
	Headers map[string]any
	Body    []byte
}

func NewRequest() Request {
	return Request{
		GetParams:  make(map[string]any),
		Headers:    make(map[string]any),
		Cookies:    make(map[string]any),
		PostParams: make(map[string]any),
	}
}

func NewResponse() Response {
	return Response{
		Headers: make(map[string]any),
		Body:    make([]byte, 0),
	}
}
