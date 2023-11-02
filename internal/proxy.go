package internal

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// ProxyServer is a very simple proxy server, just for test.
type ProxyServer struct {
	server *http.Server
}

func NewProxyServer() *ProxyServer {
	server := new(ProxyServer)

	server.server = new(http.Server)
	server.server.Addr = "127.0.0.1:8000"
	server.server.Handler = server

	return server
}

func (server *ProxyServer) Run() {
	err := server.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
}

func (server *ProxyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if err := server.validateAuth(req); err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}

	cli := &http.Client{}

	server.deleteHopHeader(req.Header)

	req.RequestURI = ""
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		server.setForwardForHeader(req.Header, host)
	}

	resp, err := cli.Do(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	server.deleteHopHeader(resp.Header)

	for k, v := range resp.Header {
		for _, v := range v {
			rw.Header().Add(k, v)
		}
	}
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

func (server *ProxyServer) validateAuth(req *http.Request) error {
	user, pass, ok := req.BasicAuth()
	if !ok {
		return nil
	}

	if user != "user" || pass != "pass" {
		return errors.New("Forbidden")
	}

	return nil
}

func (server *ProxyServer) deleteHopHeader(header http.Header) {
	for _, k := range []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"TE",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	} {
		header.Del(k)
	}
}

func (server *ProxyServer) setForwardForHeader(header http.Header, host string) {
	xff, ok := header["X-Forward-For"]
	if ok {
		host = strings.Join(xff, ", ") + ", " + host
	}
	header.Set("X-Forward-For", host)
}
