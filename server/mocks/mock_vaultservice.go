// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sordgom/PasswordManager/server/config (interfaces: VaultService)
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_vaultservice.go -package=mocks github.com/sordgom/PasswordManager/server/config VaultService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/sordgom/PasswordManager/server/model"
	gomock "go.uber.org/mock/gomock"
)

// MockVaultService is a mock of VaultService interface.
type MockVaultService struct {
	ctrl     *gomock.Controller
	recorder *MockVaultServiceMockRecorder
}

// MockVaultServiceMockRecorder is the mock recorder for MockVaultService.
type MockVaultServiceMockRecorder struct {
	mock *MockVaultService
}

// NewMockVaultService creates a new mock instance.
func NewMockVaultService(ctrl *gomock.Controller) *MockVaultService {
	mock := &MockVaultService{ctrl: ctrl}
	mock.recorder = &MockVaultServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVaultService) EXPECT() *MockVaultServiceMockRecorder {
	return m.recorder
}

// LoadVaultFromRedis mocks base method.
func (m *MockVaultService) LoadVaultFromRedis(arg0 context.Context, arg1 string) (*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadVaultFromRedis", arg0, arg1)
	ret0, _ := ret[0].(*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadVaultFromRedis indicates an expected call of LoadVaultFromRedis.
func (mr *MockVaultServiceMockRecorder) LoadVaultFromRedis(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadVaultFromRedis", reflect.TypeOf((*MockVaultService)(nil).LoadVaultFromRedis), arg0, arg1)
}

// SaveVaultToRedis mocks base method.
func (m *MockVaultService) SaveVaultToRedis(arg0 context.Context, arg1 *model.Vault) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveVaultToRedis", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveVaultToRedis indicates an expected call of SaveVaultToRedis.
func (mr *MockVaultServiceMockRecorder) SaveVaultToRedis(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveVaultToRedis", reflect.TypeOf((*MockVaultService)(nil).SaveVaultToRedis), arg0, arg1)
}
