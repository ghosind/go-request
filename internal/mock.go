package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	switch req.URL.Path {
	case "/redirect":
		server.redirectHandler(rw, req)
	case "/status":
		server.statusHandler(rw, req)
	default:
		server.defaultHandler(rw, req)
	}
}

func (server *MockServer) redirectHandler(rw http.ResponseWriter, req *http.Request) {
	tried := getIntParameter(req, "tried", 0)

	rw.Header().Set("Location", fmt.Sprintf("http://127.0.0.1:8080/redirect?tried=%d", tried+1))
	rw.WriteHeader(http.StatusFound)
}

func (server *MockServer) statusHandler(rw http.ResponseWriter, req *http.Request) {
	status := getIntParameter(req, "status", 200)

	rw.WriteHeader(int(status))
}

func (server *MockServer) defaultHandler(rw http.ResponseWriter, req *http.Request) {
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
		"userAgent":   req.Header.Get("User-Agent"),
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

func getIntParameter(req *http.Request, key string, defaultValue int64) int64 {
	queries := req.URL.Query()
	if !queries.Has(key) {
		return defaultValue
	}

	value := queries.Get(key)
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return intValue
}
