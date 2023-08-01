package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MockServer struct {
	server *http.Server
}

func NewMockServer() *MockServer {
	server := new(MockServer)

	server.server = new(http.Server)
	server.server.Addr = "127.0.0.1:8080"
	server.server.Handler = server

	return server
}

func (server *MockServer) Run() {
	err := server.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
}

func (server *MockServer) Shutdown() {
	server.server.Close()
}

func (server *MockServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}

	body := map[string]any{
		"method":      req.Method,
		"path":        req.URL.Path,
		"contentType": req.Header.Get("Content-Type"),
		"body":        string(payload),
		"query":       req.URL.RawQuery,
	}
	if token := req.Header.Get("Authorization"); token != "" {
		body["token"] = token
	}

	data, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(data); err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
}
