package proxyDelivery

import (
	log "github.com/sirupsen/logrus"
	"io"
	"main/internal/pkg/proxy"
	"net/http"
)

type ProxyServer struct {
	proxyRepo proxy.Repository
}

func (ps *ProxyServer) ListenAndServe(addr string) error {
	server := http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ps.proxyHttp(w, r)
		}),
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
	return server.ListenAndServe()
}

const ProxyHeader = "Proxy-Connection"

func copyRespToWriter(w http.ResponseWriter, resp *http.Response) error {
	header := w.Header()
	for key, values := range resp.Header {
		for _, value := range values {
			header.Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProxyServer) proxyHttp(w http.ResponseWriter, r *http.Request) {
	r.Header.Del(ProxyHeader)

	err := ps.proxyRepo.SaveRequest(r)
	if err != nil {
		log.Errorln("error saving request", err)
		return
	}

	r.RequestURI = ""
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(r)
	if err != nil {
		log.Errorln("error proxying request to another host", err)
		return
	}
	defer resp.Body.Close()

	err = ps.proxyRepo.SaveResponse(resp)
	if err != nil {
		log.Errorln("error saving response", err)
		return
	}

	err = copyRespToWriter(w, resp)
	if err != nil {
		log.Errorln("error saving response", err)
		return
	}
	return
}

func NewProxy(repository proxy.Repository) *ProxyServer {
	return &ProxyServer{proxyRepo: repository}
}
