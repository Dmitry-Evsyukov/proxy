package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"main/internal/db/init_db"
	proxyDelivery "main/internal/pkg/proxy/delivery"
	proxyRepo "main/internal/pkg/proxy/repo"
	"main/internal/pkg/scanner/delivery"
	"net/http"
)

const dbUrl = "postgres://proxy:proxy@service-db-postgres:5432/proxy"
const proxyAddr = ":8081"
const webApiAddr = ":8000"

func main() {
	pg := init_db.NewConn(dbUrl)

	pr := proxyRepo.New(pg)
	//sr := scannerRepo.New(pg)

	proxy := proxyDelivery.NewProxy(pr)

	scannerHandler := delivery.NewScannerHandler(pr)
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/requests", scannerHandler.AllRequests).Methods("GET")
		router.HandleFunc("/request/{id}", scannerHandler.GetRequest).Methods("GET")
		router.HandleFunc("/scan/{id}", scannerHandler.ScanRequest).Methods("GET")
		router.HandleFunc("/repeat/{id}", scannerHandler.RepeatRequest).Methods("GET")
		log.Fatal(http.ListenAndServe(webApiAddr, router))
	}()

	log.Fatalln(proxy.ListenAndServe(proxyAddr))
}
