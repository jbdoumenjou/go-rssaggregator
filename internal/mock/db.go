// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jbdoumenjou/go-rssaggregator/internal/database (interfaces: Querier)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination internal/mock/db.go github.com/jbdoumenjou/go-rssaggregator/internal/database Querier
//
// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	database "github.com/jbdoumenjou/go-rssaggregator/internal/database"
	gomock "go.uber.org/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// CreateFeed mocks base method.
func (m *MockQuerier) CreateFeed(arg0 context.Context, arg1 database.CreateFeedParams) (database.Feed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeed", arg0, arg1)
	ret0, _ := ret[0].(database.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFeed indicates an expected call of CreateFeed.
func (mr *MockQuerierMockRecorder) CreateFeed(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeed", reflect.TypeOf((*MockQuerier)(nil).CreateFeed), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockQuerier) CreateUser(arg0 context.Context, arg1 string) (database.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(database.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockQuerierMockRecorder) CreateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockQuerier)(nil).CreateUser), arg0, arg1)
}

// GetUserFromApiKey mocks base method.
func (m *MockQuerier) GetUserFromApiKey(arg0 context.Context, arg1 string) (database.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFromApiKey", arg0, arg1)
	ret0, _ := ret[0].(database.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFromApiKey indicates an expected call of GetUserFromApiKey.
func (mr *MockQuerierMockRecorder) GetUserFromApiKey(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFromApiKey", reflect.TypeOf((*MockQuerier)(nil).GetUserFromApiKey), arg0, arg1)
}

// GetUserFromId mocks base method.
func (m *MockQuerier) GetUserFromId(arg0 context.Context, arg1 uuid.UUID) (database.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFromId", arg0, arg1)
	ret0, _ := ret[0].(database.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFromId indicates an expected call of GetUserFromId.
func (mr *MockQuerierMockRecorder) GetUserFromId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFromId", reflect.TypeOf((*MockQuerier)(nil).GetUserFromId), arg0, arg1)
}
