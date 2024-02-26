package proxyDelivery

import (
	"bufio"
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"io"
	"main/internal/pkg/proxy"
	"main/internal/pkg/utils"
	"net/http"
	"net/http/httputil"
	"strings"
)

type ProxyServer struct {
	proxyRepo proxy.Repository
}

func (ps *ProxyServer) ListenAndServe(addr string) error {
	server := http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				ps.proxyHttps(w, r)
			} else {
				ps.proxyHttp(w, r)
			}
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

	reqId, err := ps.proxyRepo.SaveRequest(r)
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

	err = ps.proxyRepo.SaveResponse(resp, reqId)
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

func (ps *ProxyServer) proxyHttps(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	host := strings.Split(r.Host, ":")[0]

	tlsConfig, err := utils.GenTLSConf(host, r.URL.Scheme)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tlsLocalConn := tls.Server(conn, &tlsConfig)
	err = tlsLocalConn.Handshake()
	if err != nil {
		tlsLocalConn.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tlsLocalConn.Close()

	remoteConn, err := tls.Dial("tcp", r.URL.Host, &tlsConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer remoteConn.Close()

	reader := bufio.NewReader(tlsLocalConn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestByte, err := httputil.DumpRequest(request, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = remoteConn.Write(requestByte)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverReader := bufio.NewReader(remoteConn)
	response, err := http.ReadResponse(serverReader, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawResponse, err := httputil.DumpResponse(response, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tlsLocalConn.Write(rawResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.URL.Scheme = "https"
	hostAndPort := strings.Split(r.URL.Host, ":")
	request.URL.Host = hostAndPort[0]

	reqId, err := ps.proxyRepo.SaveRequest(request)
	if err != nil {
		log.Printf("Error save:  %v", err)
	}

	err = ps.proxyRepo.SaveResponse(response, reqId)
	if err != nil {
		log.Printf("Error save:  %v", err)
	}
}

func NewProxy(repository proxy.Repository) *ProxyServer {
	return &ProxyServer{proxyRepo: repository}
}
