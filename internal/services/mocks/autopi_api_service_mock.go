// Code generated by MockGen. DO NOT EDIT.
// Source: autopi_api_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	services "github.com/DIMO-Network/devices-api/internal/services"
	gomock "github.com/golang/mock/gomock"
)

// MockAutoPiAPIService is a mock of AutoPiAPIService interface.
type MockAutoPiAPIService struct {
	ctrl     *gomock.Controller
	recorder *MockAutoPiAPIServiceMockRecorder
}

// MockAutoPiAPIServiceMockRecorder is the mock recorder for MockAutoPiAPIService.
type MockAutoPiAPIServiceMockRecorder struct {
	mock *MockAutoPiAPIService
}

// NewMockAutoPiAPIService creates a new mock instance.
func NewMockAutoPiAPIService(ctrl *gomock.Controller) *MockAutoPiAPIService {
	mock := &MockAutoPiAPIService{ctrl: ctrl}
	mock.recorder = &MockAutoPiAPIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutoPiAPIService) EXPECT() *MockAutoPiAPIServiceMockRecorder {
	return m.recorder
}

// ApplyTemplate mocks base method.
func (m *MockAutoPiAPIService) ApplyTemplate(deviceID string, templateID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyTemplate", deviceID, templateID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyTemplate indicates an expected call of ApplyTemplate.
func (mr *MockAutoPiAPIServiceMockRecorder) ApplyTemplate(deviceID, templateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyTemplate", reflect.TypeOf((*MockAutoPiAPIService)(nil).ApplyTemplate), deviceID, templateID)
}

// AssociateDeviceToTemplate mocks base method.
func (m *MockAutoPiAPIService) AssociateDeviceToTemplate(deviceID string, templateID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssociateDeviceToTemplate", deviceID, templateID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssociateDeviceToTemplate indicates an expected call of AssociateDeviceToTemplate.
func (mr *MockAutoPiAPIServiceMockRecorder) AssociateDeviceToTemplate(deviceID, templateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssociateDeviceToTemplate", reflect.TypeOf((*MockAutoPiAPIService)(nil).AssociateDeviceToTemplate), deviceID, templateID)
}

// CommandRaw mocks base method.
func (m *MockAutoPiAPIService) CommandRaw(deviceID, command string, withWebhook bool) (*services.AutoPiCommandResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommandRaw", deviceID, command, withWebhook)
	ret0, _ := ret[0].(*services.AutoPiCommandResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CommandRaw indicates an expected call of CommandRaw.
func (mr *MockAutoPiAPIServiceMockRecorder) CommandRaw(deviceID, command, withWebhook interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommandRaw", reflect.TypeOf((*MockAutoPiAPIService)(nil).CommandRaw), deviceID, command, withWebhook)
}

// CommandSyncDevice mocks base method.
func (m *MockAutoPiAPIService) CommandSyncDevice(deviceID string) (*services.AutoPiCommandResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommandSyncDevice", deviceID)
	ret0, _ := ret[0].(*services.AutoPiCommandResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CommandSyncDevice indicates an expected call of CommandSyncDevice.
func (mr *MockAutoPiAPIServiceMockRecorder) CommandSyncDevice(deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommandSyncDevice", reflect.TypeOf((*MockAutoPiAPIService)(nil).CommandSyncDevice), deviceID)
}

// GetCommandStatus mocks base method.
func (m *MockAutoPiAPIService) GetCommandStatus(deviceID, jobID string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommandStatus", deviceID, jobID)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommandStatus indicates an expected call of GetCommandStatus.
func (mr *MockAutoPiAPIServiceMockRecorder) GetCommandStatus(deviceID, jobID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommandStatus", reflect.TypeOf((*MockAutoPiAPIService)(nil).GetCommandStatus), deviceID, jobID)
}

// GetDeviceByID mocks base method.
func (m *MockAutoPiAPIService) GetDeviceByID(deviceID string) (*services.AutoPiDongleDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceByID", deviceID)
	ret0, _ := ret[0].(*services.AutoPiDongleDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceByID indicates an expected call of GetDeviceByID.
func (mr *MockAutoPiAPIServiceMockRecorder) GetDeviceByID(deviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceByID", reflect.TypeOf((*MockAutoPiAPIService)(nil).GetDeviceByID), deviceID)
}

// GetDeviceByUnitID mocks base method.
func (m *MockAutoPiAPIService) GetDeviceByUnitID(unitID string) (*services.AutoPiDongleDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceByUnitID", unitID)
	ret0, _ := ret[0].(*services.AutoPiDongleDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceByUnitID indicates an expected call of GetDeviceByUnitID.
func (mr *MockAutoPiAPIServiceMockRecorder) GetDeviceByUnitID(unitID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceByUnitID", reflect.TypeOf((*MockAutoPiAPIService)(nil).GetDeviceByUnitID), unitID)
}

// PatchVehicleProfile mocks base method.
func (m *MockAutoPiAPIService) PatchVehicleProfile(vehicleID int, profile services.PatchVehicleProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchVehicleProfile", vehicleID, profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchVehicleProfile indicates an expected call of PatchVehicleProfile.
func (mr *MockAutoPiAPIServiceMockRecorder) PatchVehicleProfile(vehicleID, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchVehicleProfile", reflect.TypeOf((*MockAutoPiAPIService)(nil).PatchVehicleProfile), vehicleID, profile)
}

// UnassociateDeviceTemplate mocks base method.
func (m *MockAutoPiAPIService) UnassociateDeviceTemplate(deviceID string, templateID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnassociateDeviceTemplate", deviceID, templateID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnassociateDeviceTemplate indicates an expected call of UnassociateDeviceTemplate.
func (mr *MockAutoPiAPIServiceMockRecorder) UnassociateDeviceTemplate(deviceID, templateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnassociateDeviceTemplate", reflect.TypeOf((*MockAutoPiAPIService)(nil).UnassociateDeviceTemplate), deviceID, templateID)
}
