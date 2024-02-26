package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

//{
//"method": "POST",
//"path": "/path1/path2",
//"get_params": {
//"x": 123,
//"y": "qwe"
//},
//"headers": {
//"Host": "example.org",
//"Header": "value"
//},
//"cookies": {
//"cookie1": 1,
//"cookie2": "qwe"
//},
//"post_params": {
//"z": "zxc"
//}
//}

func RequestToJson(req *http.Request) ([]byte, error) {
	reqMap := make(map[string]any)
	reqMap["method"] = req.Method
	reqMap["path"] = req.URL.Path

	headers := make(map[string]any)
	for key, value := range req.Header {
		headers[key] = value
	}
	reqMap["headers"] = headers

	cookies := make(map[string]any)
	for index, value := range req.Cookies() {
		cookies["cookie"+strconv.Itoa(index)] = value
	}
	reqMap["cookies"] = cookies

	urlParams := make(map[string]any)
	for k, v := range req.URL.Query() {
		urlParams[k] = v
	}
	reqMap["get_params"] = urlParams

	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	postParams := make(map[string]any)
	for key, values := range req.Form {
		if len(values) > 0 {
			postParams[key] = values[0]
		}
	}
	reqMap["post_params"] = postParams

	reqJson, err := json.Marshal(reqMap)
	if err != nil {
		return nil, err
	}

	return reqJson, nil
}

//{
//"code": 200,
//"message": "OK",
//"headers": {
//"Server": "nginx/1.14.1",
//"Header": "value"
//},
//"body": "<html>..."
//}

func ResponseToJson(resp *http.Response) ([]byte, error) {
	respMap := make(map[string]any)
	respMap["code"] = resp.StatusCode
	respMap["message"] = resp.Status

	headers := make(map[string]any)
	for key, value := range resp.Header {
		headers[key] = value
	}
	respMap["headers"] = headers

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respMap["body"] = data

	respJson, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return respJson, nil
}

func JsonToResponse(data []byte) (*http.Response, error) {
	respJson := make(map[string]any)
	err := json.Unmarshal(data, &respJson)
	if err != nil {
		return nil, err
	}

	result := &http.Response{
		Status:     respJson["message"].(string),
		StatusCode: respJson["status"].(int),
	}

	for key, value := range respJson["headers"].(map[string]any) {
		result.Header.Add(key, value.(string))
	}

	buffer := bytes.NewBuffer(respJson["body"].([]byte))
	err = result.Write(buffer)
	if err != nil {
		return nil, err
	}

	return result, nil
}

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
