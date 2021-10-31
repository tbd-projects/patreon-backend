// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock_usecase_subscribers is a generated GoMock package.
package mock_usecase_subscribers

import (
	models "patreon/internal/app/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// GetCreators mocks base method.
func (m *MockUsecase) GetCreators(userID int64) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreators", userID)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCreators indicates an expected call of GetCreators.
func (mr *MockUsecaseMockRecorder) GetCreators(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreators", reflect.TypeOf((*MockUsecase)(nil).GetCreators), userID)
}

// GetSubscribers mocks base method.
func (m *MockUsecase) GetSubscribers(creatorID int64) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscribers", creatorID)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscribers indicates an expected call of GetSubscribers.
func (mr *MockUsecaseMockRecorder) GetSubscribers(creatorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribers", reflect.TypeOf((*MockUsecase)(nil).GetSubscribers), creatorID)
}

// Subscribe mocks base method.
func (m *MockUsecase) Subscribe(subscriber *models.Subscriber) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", subscriber)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockUsecaseMockRecorder) Subscribe(subscriber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockUsecase)(nil).Subscribe), subscriber)
}

// UnSubscribe mocks base method.
func (m *MockUsecase) UnSubscribe(subscriber *models.Subscriber) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnSubscribe", subscriber)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnSubscribe indicates an expected call of UnSubscribe.
func (mr *MockUsecaseMockRecorder) UnSubscribe(subscriber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnSubscribe", reflect.TypeOf((*MockUsecase)(nil).UnSubscribe), subscriber)
}