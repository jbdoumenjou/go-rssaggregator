package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

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

// This test is a poc to show how to mock a database call.
// It currently tests the CreateUser handler but should be improved.
func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()

	querier := mockdb.NewMockQuerier(ctrl)
	querier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(database.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	r := NewRouter(querier)

	user := struct {
		Name string `json:"name"`
	}{
		Name: "John Doe",
	}
	data, err := json.Marshal(user)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	var actualUser struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &actualUser)
	require.NoError(t, err)

	require.NotEmpty(t, actualUser.ID)
	require.Equal(t, user.Name, actualUser.Name)
	require.NotEmpty(t, actualUser.CreatedAt)
	require.NotEmpty(t, actualUser.UpdatedAt)
	fmt.Println(actualUser)
}

func TestUserHandler_CreateUser2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()

	querier := mockdb.NewMockQuerier(ctrl)
	querier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(database.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	r := NewRouter(querier)

	user := struct {
		Name string `json:"name"`
	}{
		Name: "John Doe",
	}
	data, err := json.Marshal(user)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	var actualUser struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &actualUser)
	require.NoError(t, err)

	require.NotEmpty(t, actualUser.ID)
	require.Equal(t, user.Name, actualUser.Name)
	require.NotEmpty(t, actualUser.CreatedAt)
	require.NotEmpty(t, actualUser.UpdatedAt)
	fmt.Println(actualUser)
}
