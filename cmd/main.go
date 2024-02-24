package main

import (
	log "github.com/sirupsen/logrus"
	proxyDelivery "main/internal/pkg/proxy/delivery"
	proxyRepo "main/internal/pkg/proxy/repo"
)

const dbUrl = "postgres://proxy:proxy@service-db-postgres:5432/proxy"
const proxyAddr = ":8081"
const webApiAddr = ":8000"

func main() {
	//pg := init_db.NewConn(dbUrl)
	pr := proxyRepo.New(nil)
	proxy := proxyDelivery.NewProxy(pr)

	//go func() {
	//	router := mux.NewRouter()
	//	router.HandleFunc("/requests", AllRequests).Methods("GET")
	//	router.HandleFunc("/request/{id}", GetRequest).Methods("GET")
	//	router.HandleFunc("/scan/{id}", ScanRequest).Methods("GET")
	//	router.HandleFunc("/repeat/{id}", RepeateRequest).Methods("GET")
	//	log.Fatal(http.ListenAndServe(webApiAddr, router))
	//}()

	log.Fatalln(proxy.ListenAndServe(proxyAddr))
}
