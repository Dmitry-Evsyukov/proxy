package delivery

import (
	"main/internal/pkg/proxy"
	"net/http"
)

type ScannerHandler struct {
	proxyRepo proxy.Repository
}

//go func() {
//	router := mux.NewRouter()
//	router.HandleFunc("/requests", AllRequests).Methods("GET")
//	router.HandleFunc("/request/{id}", GetRequest).Methods("GET")
//	router.HandleFunc("/scan/{id}", ScanRequest).Methods("GET")
//	router.HandleFunc("/repeat/{id}", RepeatRequest).Methods("GET")
//	log.Fatal(http.ListenAndServe(webApiAddr, router))
//}()

func (sh *ScannerHandler) AllRequests(w http.ResponseWriter, r *http.Request) {

}

func (sh *ScannerHandler) GetRequest(w http.ResponseWriter, r *http.Request) {

}

func (sh *ScannerHandler) ScanRequest(w http.ResponseWriter, r *http.Request) {

}

func (sh *ScannerHandler) RepeatRequest(w http.ResponseWriter, r *http.Request) {

}

func NewScannerHandler(repository proxy.Repository) *ScannerHandler {
	return &ScannerHandler{proxyRepo: repository}
}
