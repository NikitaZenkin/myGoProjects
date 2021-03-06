// Code generated by MockGen. DO NOT EDIT.
// Source: userHandler.go

// Package handlers is a generated GoMock package.
package handlers

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	user "rclone/pkg/user"
	reflect "reflect"
)

// MockUserRepositoryInterface is a mock of UserRepositoryInterface interface
type MockUserRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryInterfaceMockRecorder
}

// MockUserRepositoryInterfaceMockRecorder is the mock recorder for MockUserRepositoryInterface
type MockUserRepositoryInterfaceMockRecorder struct {
	mock *MockUserRepositoryInterface
}

// NewMockUserRepositoryInterface creates a new mock instance
func NewMockUserRepositoryInterface(ctrl *gomock.Controller) *MockUserRepositoryInterface {
	mock := &MockUserRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserRepositoryInterface) EXPECT() *MockUserRepositoryInterfaceMockRecorder {
	return m.recorder
}

// UserExist mocks base method
func (m *MockUserRepositoryInterface) UserExist(userName string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExist", userName)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserExist indicates an expected call of UserExist
func (mr *MockUserRepositoryInterfaceMockRecorder) UserExist(userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExist", reflect.TypeOf((*MockUserRepositoryInterface)(nil).UserExist), userName)
}

// AddUser mocks base method
func (m *MockUserRepositoryInterface) AddUser(newUser *user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", newUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser
func (mr *MockUserRepositoryInterfaceMockRecorder) AddUser(newUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockUserRepositoryInterface)(nil).AddUser), newUser)
}

// MockSessionsManagerInterface is a mock of SessionsManagerInterface interface
type MockSessionsManagerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsManagerInterfaceMockRecorder
}

// MockSessionsManagerInterfaceMockRecorder is the mock recorder for MockSessionsManagerInterface
type MockSessionsManagerInterfaceMockRecorder struct {
	mock *MockSessionsManagerInterface
}

// NewMockSessionsManagerInterface creates a new mock instance
func NewMockSessionsManagerInterface(ctrl *gomock.Controller) *MockSessionsManagerInterface {
	mock := &MockSessionsManagerInterface{ctrl: ctrl}
	mock.recorder = &MockSessionsManagerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSessionsManagerInterface) EXPECT() *MockSessionsManagerInterfaceMockRecorder {
	return m.recorder
}

// CreateSession mocks base method
func (m *MockSessionsManagerInterface) CreateSession(w http.ResponseWriter, userID, userName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", w, userID, userName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession
func (mr *MockSessionsManagerInterfaceMockRecorder) CreateSession(w, userID, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionsManagerInterface)(nil).CreateSession), w, userID, userName)
}
