package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var mux *http.ServeMux
var writer *httptest.ResponseRecorder

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
}

func TestGetHealthCheck(t *testing.T) {
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/healthcheck", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.Bytes()
	if string(body) != "ok" {
		t.Errorf("Response body not ok")
	}
}

func TestPostHealthCheck(t *testing.T) {
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/healthcheck", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.Bytes()
	if string(body) != "ok" {
		t.Errorf("Response body not ok '%s'", string(body))
	}
}
