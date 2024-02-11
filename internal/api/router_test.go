package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
	assert.Equal(t, expected, rr.Body.String())
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
	assert.Equal(t, expected, rr.Body.String())
}

type user struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func createUser(t *testing.T, r http.Handler) user {
	t.Helper()

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
	require.Equal(t, http.StatusOK, rr.Code)

	var actualUser struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		ApiKey    string `json:"api_key"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &actualUser)
	require.NoError(t, err)

	require.NotEmpty(t, actualUser.ID)
	require.Equal(t, user.Name, actualUser.Name)
	require.NotEmpty(t, actualUser.CreatedAt)
	require.NotEmpty(t, actualUser.UpdatedAt)

	return actualUser
}

func TestUserHandler_db_CreateUser(t *testing.T) {
	createUser(t, NewRouter(testQueries))
}

// This test is a poc to show how to mock a database call.
// It currently tests the CreateUser handler but should be improved.
func TestUserHandler_mock_CreateUser(t *testing.T) {
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

	createUser(t, NewRouter(querier))
}

func TestUserHandler_GetUser(t *testing.T) {
	r := NewRouter(testQueries)
	user1 := createUser(t, r)

	tests := []struct {
		name               string
		setHeader          func(h http.Header)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "valid",
			setHeader: func(h http.Header) {
				h.Set("Authorization", "ApiKey "+user1.ApiKey)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: func() string {
				b, err := json.Marshal(user1)
				require.NoError(t, err)
				return string(b)
			}(),
		},
		{
			name:               "without authorization header",
			setHeader:          func(h http.Header) {},
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"error":"Forbidden"}`,
		},
		{
			name: "with empty authorization value",
			setHeader: func(h http.Header) {
				h.Set("Authorization", "")
			},
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"error":"Forbidden"}`,
		},
		{
			name: "with bad authorization key",
			setHeader: func(h http.Header) {
				h.Set("Authorization", "token "+user1.ApiKey)
			},
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"error":"Forbidden"}`,
		},
		{
			name: "with bad authorization value",
			setHeader: func(h http.Header) {
				h.Set("Authorization", "ApiKey unknown")
			},
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"error":"Forbidden"}`,
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/v1/users", http.NoBody)
			if err != nil {
				t.Fatal(err)
			}

			tc.setHeader(req.Header)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)
			assert.Equal(t, tc.expectedStatusCode, rr.Code)
			assert.JSONEq(t, tc.expectedBody, rr.Body.String())
		})
	}
}
