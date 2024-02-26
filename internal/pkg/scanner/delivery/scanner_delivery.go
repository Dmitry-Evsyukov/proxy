package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"main/internal/pkg/scanner"
	"main/internal/pkg/utils"
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
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
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
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := sh.scannerRepo.GetRequest(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	example, err := sh.scannerRepo.GetResponse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("error repeating request", err)
		return
	}

	requestFirst, err := utils.StructToRequest(req, `"`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("error repeating request", err)
		return
	}

	response, err := client.Do(requestFirst)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("error repeating request", err)
		return
	}
	if response.StatusCode != example.Code || response.Status != example.Message {
		w.WriteHeader(205)
		return
	}

	requestSecond, err := utils.StructToRequest(req, `'`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("error repeating request", err)
		return
	}

	response, err = client.Do(requestSecond)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("error repeating request", err)
		return
	}
	if response.StatusCode != example.Code || response.Status != example.Message {
		w.WriteHeader(205)
		return
	}

	return
}

func (sh *ScannerHandler) RepeatRequest(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := sh.scannerRepo.GetRequest(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := utils.StructToRequest(req, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("error repeating request", err)
		return
	}

	result, err := utils.ResponseToStruct(resp, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln("error repeating request", err)
		return
	}

	data, err := json.Marshal(result)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func NewScannerHandler(repository scanner.Repository) *ScannerHandler {
	return &ScannerHandler{scannerRepo: repository}
}
