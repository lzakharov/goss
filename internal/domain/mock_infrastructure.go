// Code generated by MockGen. DO NOT EDIT.
// Source: infrastructure.go

// Package domain is a generated GoMock package.
package domain

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMortal is a mock of Mortal interface
type MockMortal struct {
	ctrl     *gomock.Controller
	recorder *MockMortalMockRecorder
}

// MockMortalMockRecorder is the mock recorder for MockMortal
type MockMortalMockRecorder struct {
	mock *MockMortal
}

// NewMockMortal creates a new mock instance
func NewMockMortal(ctrl *gomock.Controller) *MockMortal {
	mock := &MockMortal{ctrl: ctrl}
	mock.recorder = &MockMortalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMortal) EXPECT() *MockMortalMockRecorder {
	return m.recorder
}

// IsAlive mocks base method
func (m *MockMortal) IsAlive() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAlive")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAlive indicates an expected call of IsAlive
func (mr *MockMortalMockRecorder) IsAlive() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAlive", reflect.TypeOf((*MockMortal)(nil).IsAlive))
}

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// IsAlive mocks base method
func (m *MockStorage) IsAlive() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAlive")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAlive indicates an expected call of IsAlive
func (mr *MockStorageMockRecorder) IsAlive() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAlive", reflect.TypeOf((*MockStorage)(nil).IsAlive))
}

// GetUser mocks base method
func (m *MockStorage) GetUser(userID int64) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", userID)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockStorageMockRecorder) GetUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStorage)(nil).GetUser), userID)
}

// GetUserByCredentials mocks base method
func (m *MockStorage) GetUserByCredentials(credentials *Credentials) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByCredentials", credentials)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByCredentials indicates an expected call of GetUserByCredentials
func (mr *MockStorageMockRecorder) GetUserByCredentials(credentials interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByCredentials", reflect.TypeOf((*MockStorage)(nil).GetUserByCredentials), credentials)
}

// MockSecurity is a mock of Security interface
type MockSecurity struct {
	ctrl     *gomock.Controller
	recorder *MockSecurityMockRecorder
}

// MockSecurityMockRecorder is the mock recorder for MockSecurity
type MockSecurityMockRecorder struct {
	mock *MockSecurity
}

// NewMockSecurity creates a new mock instance
func NewMockSecurity(ctrl *gomock.Controller) *MockSecurity {
	mock := &MockSecurity{ctrl: ctrl}
	mock.recorder = &MockSecurityMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecurity) EXPECT() *MockSecurityMockRecorder {
	return m.recorder
}

// IsAlive mocks base method
func (m *MockSecurity) IsAlive() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAlive")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAlive indicates an expected call of IsAlive
func (mr *MockSecurityMockRecorder) IsAlive() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAlive", reflect.TypeOf((*MockSecurity)(nil).IsAlive))
}

// CreateAuthData mocks base method
func (m *MockSecurity) CreateAuthData(user *User) (*AuthData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAuthData", user)
	ret0, _ := ret[0].(*AuthData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAuthData indicates an expected call of CreateAuthData
func (mr *MockSecurityMockRecorder) CreateAuthData(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAuthData", reflect.TypeOf((*MockSecurity)(nil).CreateAuthData), user)
}

// GetAccessTokenClaims mocks base method
func (m *MockSecurity) GetAccessTokenClaims(accessToken string) (*AccessTokenClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessTokenClaims", accessToken)
	ret0, _ := ret[0].(*AccessTokenClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccessTokenClaims indicates an expected call of GetAccessTokenClaims
func (mr *MockSecurityMockRecorder) GetAccessTokenClaims(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessTokenClaims", reflect.TypeOf((*MockSecurity)(nil).GetAccessTokenClaims), accessToken)
}

// GetRefreshTokenClaims mocks base method
func (m *MockSecurity) GetRefreshTokenClaims(refreshToken string) (*RefreshTokenClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshTokenClaims", refreshToken)
	ret0, _ := ret[0].(*RefreshTokenClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefreshTokenClaims indicates an expected call of GetRefreshTokenClaims
func (mr *MockSecurityMockRecorder) GetRefreshTokenClaims(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshTokenClaims", reflect.TypeOf((*MockSecurity)(nil).GetRefreshTokenClaims), refreshToken)
}

// InvalidateUserAuthData mocks base method
func (m *MockSecurity) InvalidateUserAuthData(userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateUserAuthData", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateUserAuthData indicates an expected call of InvalidateUserAuthData
func (mr *MockSecurityMockRecorder) InvalidateUserAuthData(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateUserAuthData", reflect.TypeOf((*MockSecurity)(nil).InvalidateUserAuthData), userID)
}
