// Code generated by MockGen. DO NOT EDIT.
// Source: device_definitions_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	models "github.com/DIMO-INC/devices-api/models"
	gomock "github.com/golang/mock/gomock"
	boil "github.com/volatiletech/sqlboiler/v4/boil"
)

// MockIDeviceDefinitionService is a mock of IDeviceDefinitionService interface.
type MockIDeviceDefinitionService struct {
	ctrl     *gomock.Controller
	recorder *MockIDeviceDefinitionServiceMockRecorder
}

// MockIDeviceDefinitionServiceMockRecorder is the mock recorder for MockIDeviceDefinitionService.
type MockIDeviceDefinitionServiceMockRecorder struct {
	mock *MockIDeviceDefinitionService
}

// NewMockIDeviceDefinitionService creates a new mock instance.
func NewMockIDeviceDefinitionService(ctrl *gomock.Controller) *MockIDeviceDefinitionService {
	mock := &MockIDeviceDefinitionService{ctrl: ctrl}
	mock.recorder = &MockIDeviceDefinitionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDeviceDefinitionService) EXPECT() *MockIDeviceDefinitionServiceMockRecorder {
	return m.recorder
}

// CheckAndSetImage mocks base method.
func (m *MockIDeviceDefinitionService) CheckAndSetImage(dd *models.DeviceDefinition) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAndSetImage", dd)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckAndSetImage indicates an expected call of CheckAndSetImage.
func (mr *MockIDeviceDefinitionServiceMockRecorder) CheckAndSetImage(dd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAndSetImage", reflect.TypeOf((*MockIDeviceDefinitionService)(nil).CheckAndSetImage), dd)
}

// FindDeviceDefinitionByMMY mocks base method.
func (m *MockIDeviceDefinitionService) FindDeviceDefinitionByMMY(ctx context.Context, db boil.ContextExecutor, mk, model string, year int, loadIntegrations bool) (*models.DeviceDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeviceDefinitionByMMY", ctx, db, mk, model, year, loadIntegrations)
	ret0, _ := ret[0].(*models.DeviceDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeviceDefinitionByMMY indicates an expected call of FindDeviceDefinitionByMMY.
func (mr *MockIDeviceDefinitionServiceMockRecorder) FindDeviceDefinitionByMMY(ctx, db, mk, model, year, loadIntegrations interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeviceDefinitionByMMY", reflect.TypeOf((*MockIDeviceDefinitionService)(nil).FindDeviceDefinitionByMMY), ctx, db, mk, model, year, loadIntegrations)
}

// UpdateDeviceDefinitionFromNHTSA mocks base method.
func (m *MockIDeviceDefinitionService) UpdateDeviceDefinitionFromNHTSA(ctx context.Context, deviceDefinitionId, vin string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDeviceDefinitionFromNHTSA", ctx, deviceDefinitionId, vin)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDeviceDefinitionFromNHTSA indicates an expected call of UpdateDeviceDefinitionFromNHTSA.
func (mr *MockIDeviceDefinitionServiceMockRecorder) UpdateDeviceDefinitionFromNHTSA(ctx, deviceDefinitionId, vin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeviceDefinitionFromNHTSA", reflect.TypeOf((*MockIDeviceDefinitionService)(nil).UpdateDeviceDefinitionFromNHTSA), ctx, deviceDefinitionId, vin)
}
