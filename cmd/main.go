package main

import (
	log "github.com/sirupsen/logrus"
	proxyDelivery "main/internal/pkg/proxy/delivery"
	proxyRepo "main/internal/pkg/proxy/repo"
)

const dbUrl = ""
const proxyAddr = ":8080"
const webApiAddr = ":8000"

func main() {
	//pg := init_db.NewConn(dbUrl)
	pr := proxyRepo.New(nil)
	proxy := proxyDelivery.NewProxy(pr)
	log.Fatalln(proxy.ListenAndServe(proxyAddr))
}
