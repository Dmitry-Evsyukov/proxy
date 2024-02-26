package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"main/internal/models"
	"net/http"
	"strconv"
)

func RequestToStruct(req *http.Request) (models.Request, error) {
	result := models.NewRequest()

	result.Method = req.Method
	result.Url = req.URL.String()
	result.Scheme = req.Proto

	for key, value := range req.Header {
		result.Headers[key] = value
	}

	for index, value := range req.Cookies() {
		result.Cookies["cookie"+strconv.Itoa(index)] = value
	}

	for k, v := range req.URL.Query() {
		result.GetParams[k] = v
	}

	err := req.ParseForm()
	if err != nil {
		return models.Request{}, err
	}

	for key, values := range req.Form {
		if len(values) > 0 {
			result.PostParams[key] = values[0]
		}
	}

	return result, nil
}

func ResponseToStruct(resp *http.Response, reqId int) (models.Response, error) {
	result := models.NewResponse()
	result.Code = resp.StatusCode
	result.Id = reqId
	result.Message = resp.Status

	for key, value := range resp.Header {
		result.Headers[key] = value
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}
	result.Body = data

	resp.Body = io.NopCloser(bytes.NewBuffer(data))

	return result, nil
}

func StructToRequest(req models.Request, symbol string) (*http.Request, error) {
	params, err := json.Marshal(req.PostParams)

	result, err := http.NewRequest(req.Method, req.Url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		if val, ok := v.(string); ok {
			result.Header.Add(k, val+symbol)
		}
	}

	for k, v := range req.Cookies {
		result.AddCookie(&http.Cookie{
			Name:  k,
			Value: v.(string) + symbol,
		})
	}

	return result, nil
}
