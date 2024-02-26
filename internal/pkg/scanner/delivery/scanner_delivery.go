package delivery

import (
	"main/internal/pkg/scanner"
	"net/http"
)

type ScannerHandler struct {
	scannerRepo scanner.Repository
}

func (sh *ScannerHandler) AllRequests(w http.ResponseWriter, r *http.Request) {
	ids, err := sh.scannerRepo.GetAllIds()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, id := range ids {

	}
}

func (sh *ScannerHandler) GetRequest(w http.ResponseWriter, r *http.Request) {

}

func (sh *ScannerHandler) ScanRequest(w http.ResponseWriter, r *http.Request) {

}

func (sh *ScannerHandler) RepeatRequest(w http.ResponseWriter, r *http.Request) {

}

func NewScannerHandler(repository scanner.Repository) *ScannerHandler {
	return &ScannerHandler{scannerRepo: repository}
}
