// Code generated by MockGen. DO NOT EDIT.
// Source: patreon/internal/app/usecase/usecase_factory (interfaces: RepositoryFactory)

// Package mock_repository_factory is a generated GoMock package.
package mock_repository_factory

import (
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	repository_access "patreon/internal/app/repository/access"
	repository_awards "patreon/internal/app/repository/awards"
	repository_creator "patreon/internal/app/repository/creator"
	repository_files "patreon/internal/app/repository/files"
	repository_info "patreon/internal/app/repository/info"
	repository_likes "patreon/internal/app/repository/likes"
	repository_payments "patreon/internal/app/repository/payments"
	repository_posts "patreon/internal/app/repository/posts"
	repository_posts_data "patreon/internal/app/repository/posts_data"
	repository_subscribers "patreon/internal/app/repository/subscribers"
	repository_user "patreon/internal/app/repository/user"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryFactory is a mock of RepositoryFactory interface.
type MockRepositoryFactory struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryFactoryMockRecorder
}

// MockRepositoryFactoryMockRecorder is the mock recorder for MockRepositoryFactory.
type MockRepositoryFactoryMockRecorder struct {
	mock *MockRepositoryFactory
}

// NewMockRepositoryFactory creates a new mock instance.
func NewMockRepositoryFactory(ctrl *gomock.Controller) *MockRepositoryFactory {
	mock := &MockRepositoryFactory{ctrl: ctrl}
	mock.recorder = &MockRepositoryFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryFactory) EXPECT() *MockRepositoryFactoryMockRecorder {
	return m.recorder
}

// GetAccessRepository mocks base method.
func (m *MockRepositoryFactory) GetAccessRepository() repository_access.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessRepository")
	ret0, _ := ret[0].(repository_access.Repository)
	return ret0
}

// GetAccessRepository indicates an expected call of GetAccessRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetAccessRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetAccessRepository))
}

// GetAwardsRepository mocks base method.
func (m *MockRepositoryFactory) GetAwardsRepository() repository_awards.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAwardsRepository")
	ret0, _ := ret[0].(repository_awards.Repository)
	return ret0
}

// GetAwardsRepository indicates an expected call of GetAwardsRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetAwardsRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAwardsRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetAwardsRepository))
}

// GetCreatorRepository mocks base method.
func (m *MockRepositoryFactory) GetCreatorRepository() repository_creator.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreatorRepository")
	ret0, _ := ret[0].(repository_creator.Repository)
	return ret0
}

// GetCreatorRepository indicates an expected call of GetCreatorRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetCreatorRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreatorRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetCreatorRepository))
}

// GetCsrfRepository mocks base method.
func (m *MockRepositoryFactory) GetCsrfRepository() repository_jwt.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCsrfRepository")
	ret0, _ := ret[0].(repository_jwt.Repository)
	return ret0
}

// GetCsrfRepository indicates an expected call of GetCsrfRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetCsrfRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCsrfRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetCsrfRepository))
}

// GetFilesRepository mocks base method.
func (m *MockRepositoryFactory) GetFilesRepository() repository_files.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilesRepository")
	ret0, _ := ret[0].(repository_files.Repository)
	return ret0
}

// GetFilesRepository indicates an expected call of GetFilesRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetFilesRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilesRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetFilesRepository))
}

// GetInfoRepository mocks base method.
func (m *MockRepositoryFactory) GetInfoRepository() repository_info.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfoRepository")
	ret0, _ := ret[0].(repository_info.Repository)
	return ret0
}

// GetInfoRepository indicates an expected call of GetInfoRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetInfoRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfoRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetInfoRepository))
}

// GetLikesRepository mocks base method.
func (m *MockRepositoryFactory) GetLikesRepository() repository_likes.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikesRepository")
	ret0, _ := ret[0].(repository_likes.Repository)
	return ret0
}

// GetLikesRepository indicates an expected call of GetLikesRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetLikesRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikesRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetLikesRepository))
}

// GetPaymentsRepository mocks base method.
func (m *MockRepositoryFactory) GetPaymentsRepository() repository_payments.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentsRepository")
	ret0, _ := ret[0].(repository_payments.Repository)
	return ret0
}

// GetPaymentsRepository indicates an expected call of GetPaymentsRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetPaymentsRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentsRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetPaymentsRepository))
}

// GetPostsDataRepository mocks base method.
func (m *MockRepositoryFactory) GetPostsDataRepository() repository_posts_data.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostsDataRepository")
	ret0, _ := ret[0].(repository_posts_data.Repository)
	return ret0
}

// GetPostsDataRepository indicates an expected call of GetPostsDataRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetPostsDataRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostsDataRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetPostsDataRepository))
}

// GetPostsRepository mocks base method.
func (m *MockRepositoryFactory) GetPostsRepository() repository_posts.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostsRepository")
	ret0, _ := ret[0].(repository_posts.Repository)
	return ret0
}

// GetPostsRepository indicates an expected call of GetPostsRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetPostsRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostsRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetPostsRepository))
}

// GetSubscribersRepository mocks base method.
func (m *MockRepositoryFactory) GetSubscribersRepository() repository_subscribers.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscribersRepository")
	ret0, _ := ret[0].(repository_subscribers.Repository)
	return ret0
}

// GetSubscribersRepository indicates an expected call of GetSubscribersRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetSubscribersRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribersRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetSubscribersRepository))
}

// GetUserRepository mocks base method.
func (m *MockRepositoryFactory) GetUserRepository() repository_user.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserRepository")
	ret0, _ := ret[0].(repository_user.Repository)
	return ret0
}

// GetUserRepository indicates an expected call of GetUserRepository.
func (mr *MockRepositoryFactoryMockRecorder) GetUserRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserRepository", reflect.TypeOf((*MockRepositoryFactory)(nil).GetUserRepository))
}
