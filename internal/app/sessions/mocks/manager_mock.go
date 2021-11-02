// Code generated by MockGen. DO NOT EDIT.
// Source: patreon/internal/app/sessions (interfaces: SessionsManager)

// Package mocks is a generated GoMock package.
package mock_sessions

import (
	models "patreon/internal/app/sessions/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionsManager is a mock of SessionsManager interface.
type MockSessionsManager struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsManagerMockRecorder
}

// MockSessionsManagerMockRecorder is the mock recorder for MockSessionsManager.
type MockSessionsManagerMockRecorder struct {
	mock *MockSessionsManager
}

// NewMockSessionsManager creates a new mock instance.
func NewMockSessionsManager(ctrl *gomock.Controller) *MockSessionsManager {
	mock := &MockSessionsManager{ctrl: ctrl}
	mock.recorder = &MockSessionsManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionsManager) EXPECT() *MockSessionsManagerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockSessionsManager) Check(arg0 string) (models.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0)
	ret0, _ := ret[0].(models.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockSessionsManagerMockRecorder) Check(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockSessionsManager)(nil).Check), arg0)
}

// Create mocks base method.
func (m *MockSessionsManager) Create(arg0 int64) (models.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(models.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionsManagerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionsManager)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockSessionsManager) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionsManagerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessionsManager)(nil).Delete), arg0)
}
