package delivery

import (
	"encoding/json"
	"main/internal/pkg/scanner"
	"net/http"
	"strconv"
)

type ScannerHandler struct {
	scannerRepo scanner.Repository
}

func (sh *ScannerHandler) AllRequests(w http.ResponseWriter, r *http.Request) {
	reqs, err := sh.scannerRepo.GetAllRequests()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(reqs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (sh *ScannerHandler) GetRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req, err := sh.scannerRepo.GetRequest(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (sh *ScannerHandler) ScanRequest(w http.ResponseWriter, r *http.Request) {

}

func (sh *ScannerHandler) RepeatRequest(w http.ResponseWriter, r *http.Request) {

}

func NewScannerHandler(repository scanner.Repository) *ScannerHandler {
	return &ScannerHandler{scannerRepo: repository}
}
