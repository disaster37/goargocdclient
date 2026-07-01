package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-resty/resty/v2"
)

func newTestClient(handler http.Handler) (*resty.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := resty.New().
		SetBaseURL(server.URL).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")
	return client, server
}

func newFailingClient() *resty.Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	client := resty.New().
		SetBaseURL(server.URL).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")
	server.Close()
	return client
}

func jsonResponse(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		json.NewEncoder(w).Encode(body)
	}
}
