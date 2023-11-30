package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jbdoumenjou/go-rssaggregator/internal/database"

	mockdb "github.com/jbdoumenjou/go-rssaggregator/internal/mock"
	"go.uber.org/mock/gomock"
)

func TestReadiness(t *testing.T) {
	r := NewRouter(nil)
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
	r := NewRouter(nil)
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

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querier := mockdb.NewMockQuerier(ctrl)
	querier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(database.User{}, nil)
}
