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

//func JsonToResponse(data []byte) (*http.Response, error) {
//	respJson := make(map[string]any)
//	err := json.Unmarshal(data, &respJson)
//	if err != nil {
//		return nil, err
//	}
//
//	result := &http.Response{
//		Status:     respJson["message"].(string),
//		StatusCode: respJson["status"].(int),
//	}
//
//	for key, value := range respJson["headers"].(map[string]any) {
//		result.Header.Add(key, value.(string))
//	}
//
//	buffer := bytes.NewBuffer(respJson["body"].([]byte))
//	err = result.Write(buffer)
//	if err != nil {
//		return nil, err
//	}
//
//	return result, nil
//}

func JsonToRequest(data []byte) (*http.Request, error) {
	reqJson := make(map[string]any)
	err := json.Unmarshal(data, &reqJson)
	if err != nil {
		return nil, err
	}

	result, err := http.NewRequest(reqJson["method"].(string), reqJson["path"].(string), bytes.NewReader(reqJson["post_params"].([]byte)))
	if err != nil {
		return nil, err
	}
	return result, nil
}
