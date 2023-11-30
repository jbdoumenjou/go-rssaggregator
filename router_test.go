package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadiness(t *testing.T) {
	r := NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/v1/readiness", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("Expected body to contain 'ok', got %s", rr.Body.String())
	}
}

func TestErr(t *testing.T) {
	r := NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/v1/err", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", status)
	}

	expected := `{"error":"Internal Server Error"}`
	if rr.Body.String() != expected {
		t.Errorf("Expected body to contain 'Internal Server Error', got %s", rr.Body.String())
	}
}
