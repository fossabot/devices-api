// Code generated by MockGen. DO NOT EDIT.
// Source: tesla_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	services "github.com/DIMO-Network/devices-api/internal/services"
	gomock "github.com/golang/mock/gomock"
)

// MockTeslaService is a mock of TeslaService interface.
type MockTeslaService struct {
	ctrl     *gomock.Controller
	recorder *MockTeslaServiceMockRecorder
}

// MockTeslaServiceMockRecorder is the mock recorder for MockTeslaService.
type MockTeslaServiceMockRecorder struct {
	mock *MockTeslaService
}

// NewMockTeslaService creates a new mock instance.
func NewMockTeslaService(ctrl *gomock.Controller) *MockTeslaService {
	mock := &MockTeslaService{ctrl: ctrl}
	mock.recorder = &MockTeslaServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTeslaService) EXPECT() *MockTeslaServiceMockRecorder {
	return m.recorder
}

// GetVehicle mocks base method.
func (m *MockTeslaService) GetVehicle(ownerAccessToken string, id int) (*services.TeslaVehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVehicle", ownerAccessToken, id)
	ret0, _ := ret[0].(*services.TeslaVehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVehicle indicates an expected call of GetVehicle.
func (mr *MockTeslaServiceMockRecorder) GetVehicle(ownerAccessToken, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVehicle", reflect.TypeOf((*MockTeslaService)(nil).GetVehicle), ownerAccessToken, id)
}

// LockDoor mocks base method.
func (m *MockTeslaService) LockDoor(ownerAccessToken string, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockDoor", ownerAccessToken, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// LockDoor indicates an expected call of LockDoor.
func (mr *MockTeslaServiceMockRecorder) LockDoor(ownerAccessToken, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockDoor", reflect.TypeOf((*MockTeslaService)(nil).LockDoor), ownerAccessToken, id)
}

// UnlockDoor mocks base method.
func (m *MockTeslaService) UnlockDoor(ownerAccessToken string, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlockDoor", ownerAccessToken, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlockDoor indicates an expected call of UnlockDoor.
func (mr *MockTeslaServiceMockRecorder) UnlockDoor(ownerAccessToken, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlockDoor", reflect.TypeOf((*MockTeslaService)(nil).UnlockDoor), ownerAccessToken, id)
}

// WakeUpVehicle mocks base method.
func (m *MockTeslaService) WakeUpVehicle(ownerAccessToken string, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WakeUpVehicle", ownerAccessToken, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// WakeUpVehicle indicates an expected call of WakeUpVehicle.
func (mr *MockTeslaServiceMockRecorder) WakeUpVehicle(ownerAccessToken, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WakeUpVehicle", reflect.TypeOf((*MockTeslaService)(nil).WakeUpVehicle), ownerAccessToken, id)
}
