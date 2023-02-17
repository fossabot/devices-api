// Code generated by MockGen. DO NOT EDIT.
// Source: device_definitions_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	big "math/big"
	reflect "reflect"

	grpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	services "github.com/DIMO-Network/devices-api/internal/services"
	gomock "github.com/golang/mock/gomock"
	boil "github.com/volatiletech/sqlboiler/v4/boil"
)

// MockDeviceDefinitionService is a mock of DeviceDefinitionService interface.
type MockDeviceDefinitionService struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceDefinitionServiceMockRecorder
}

// MockDeviceDefinitionServiceMockRecorder is the mock recorder for MockDeviceDefinitionService.
type MockDeviceDefinitionServiceMockRecorder struct {
	mock *MockDeviceDefinitionService
}

// NewMockDeviceDefinitionService creates a new mock instance.
func NewMockDeviceDefinitionService(ctrl *gomock.Controller) *MockDeviceDefinitionService {
	mock := &MockDeviceDefinitionService{ctrl: ctrl}
	mock.recorder = &MockDeviceDefinitionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceDefinitionService) EXPECT() *MockDeviceDefinitionServiceMockRecorder {
	return m.recorder
}

// CreateIntegration mocks base method.
func (m *MockDeviceDefinitionService) CreateIntegration(ctx context.Context, integrationType, vendor, style string) (*grpc.Integration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIntegration", ctx, integrationType, vendor, style)
	ret0, _ := ret[0].(*grpc.Integration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIntegration indicates an expected call of CreateIntegration.
func (mr *MockDeviceDefinitionServiceMockRecorder) CreateIntegration(ctx, integrationType, vendor, style interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIntegration", reflect.TypeOf((*MockDeviceDefinitionService)(nil).CreateIntegration), ctx, integrationType, vendor, style)
}

// DecodeVIN mocks base method.
func (m *MockDeviceDefinitionService) DecodeVIN(ctx context.Context, vin string) (*grpc.DecodeVinResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeVIN", ctx, vin)
	ret0, _ := ret[0].(*grpc.DecodeVinResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecodeVIN indicates an expected call of DecodeVIN.
func (mr *MockDeviceDefinitionServiceMockRecorder) DecodeVIN(ctx, vin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeVIN", reflect.TypeOf((*MockDeviceDefinitionService)(nil).DecodeVIN), ctx, vin)
}

// FindDeviceDefinitionByMMY mocks base method.
func (m *MockDeviceDefinitionService) FindDeviceDefinitionByMMY(ctx context.Context, mk, model string, year int) (*grpc.GetDeviceDefinitionItemResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeviceDefinitionByMMY", ctx, mk, model, year)
	ret0, _ := ret[0].(*grpc.GetDeviceDefinitionItemResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeviceDefinitionByMMY indicates an expected call of FindDeviceDefinitionByMMY.
func (mr *MockDeviceDefinitionServiceMockRecorder) FindDeviceDefinitionByMMY(ctx, mk, model, year interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeviceDefinitionByMMY", reflect.TypeOf((*MockDeviceDefinitionService)(nil).FindDeviceDefinitionByMMY), ctx, mk, model, year)
}

// GetDeviceDefinitionByID mocks base method.
func (m *MockDeviceDefinitionService) GetDeviceDefinitionByID(ctx context.Context, id string) (*grpc.GetDeviceDefinitionItemResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceDefinitionByID", ctx, id)
	ret0, _ := ret[0].(*grpc.GetDeviceDefinitionItemResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceDefinitionByID indicates an expected call of GetDeviceDefinitionByID.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetDeviceDefinitionByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceDefinitionByID", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetDeviceDefinitionByID), ctx, id)
}

// GetDeviceDefinitionsByIDs mocks base method.
func (m *MockDeviceDefinitionService) GetDeviceDefinitionsByIDs(ctx context.Context, ids []string) ([]*grpc.GetDeviceDefinitionItemResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceDefinitionsByIDs", ctx, ids)
	ret0, _ := ret[0].([]*grpc.GetDeviceDefinitionItemResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceDefinitionsByIDs indicates an expected call of GetDeviceDefinitionsByIDs.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetDeviceDefinitionsByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceDefinitionsByIDs", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetDeviceDefinitionsByIDs), ctx, ids)
}

// GetIntegrationByFilter mocks base method.
func (m *MockDeviceDefinitionService) GetIntegrationByFilter(ctx context.Context, integrationType, vendor, style string) (*grpc.Integration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntegrationByFilter", ctx, integrationType, vendor, style)
	ret0, _ := ret[0].(*grpc.Integration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIntegrationByFilter indicates an expected call of GetIntegrationByFilter.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetIntegrationByFilter(ctx, integrationType, vendor, style interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntegrationByFilter", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetIntegrationByFilter), ctx, integrationType, vendor, style)
}

// GetIntegrationByID mocks base method.
func (m *MockDeviceDefinitionService) GetIntegrationByID(ctx context.Context, id string) (*grpc.Integration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntegrationByID", ctx, id)
	ret0, _ := ret[0].(*grpc.Integration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIntegrationByID indicates an expected call of GetIntegrationByID.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetIntegrationByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntegrationByID", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetIntegrationByID), ctx, id)
}

// GetIntegrationByVendor mocks base method.
func (m *MockDeviceDefinitionService) GetIntegrationByVendor(ctx context.Context, vendor string) (*grpc.Integration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntegrationByVendor", ctx, vendor)
	ret0, _ := ret[0].(*grpc.Integration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIntegrationByVendor indicates an expected call of GetIntegrationByVendor.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetIntegrationByVendor(ctx, vendor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntegrationByVendor", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetIntegrationByVendor), ctx, vendor)
}

// GetIntegrations mocks base method.
func (m *MockDeviceDefinitionService) GetIntegrations(ctx context.Context) ([]*grpc.Integration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntegrations", ctx)
	ret0, _ := ret[0].([]*grpc.Integration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIntegrations indicates an expected call of GetIntegrations.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetIntegrations(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntegrations", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetIntegrations), ctx)
}

// GetMakeByTokenID mocks base method.
func (m *MockDeviceDefinitionService) GetMakeByTokenID(ctx context.Context, tokenID *big.Int) (*grpc.DeviceMake, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMakeByTokenID", ctx, tokenID)
	ret0, _ := ret[0].(*grpc.DeviceMake)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMakeByTokenID indicates an expected call of GetMakeByTokenID.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetMakeByTokenID(ctx, tokenID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMakeByTokenID", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetMakeByTokenID), ctx, tokenID)
}

// GetOrCreateMake mocks base method.
func (m *MockDeviceDefinitionService) GetOrCreateMake(ctx context.Context, tx boil.ContextExecutor, makeName string) (*grpc.DeviceMake, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrCreateMake", ctx, tx, makeName)
	ret0, _ := ret[0].(*grpc.DeviceMake)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrCreateMake indicates an expected call of GetOrCreateMake.
func (mr *MockDeviceDefinitionServiceMockRecorder) GetOrCreateMake(ctx, tx, makeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrCreateMake", reflect.TypeOf((*MockDeviceDefinitionService)(nil).GetOrCreateMake), ctx, tx, makeName)
}

// PullDrivlyData mocks base method.
func (m *MockDeviceDefinitionService) PullDrivlyData(ctx context.Context, userDeviceID, deviceDefinitionID, vin string, forceSetAll bool) (services.DrivlyDataStatusEnum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullDrivlyData", ctx, userDeviceID, deviceDefinitionID, vin, forceSetAll)
	ret0, _ := ret[0].(services.DrivlyDataStatusEnum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PullDrivlyData indicates an expected call of PullDrivlyData.
func (mr *MockDeviceDefinitionServiceMockRecorder) PullDrivlyData(ctx, userDeviceID, deviceDefinitionID, vin, forceSetAll interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullDrivlyData", reflect.TypeOf((*MockDeviceDefinitionService)(nil).PullDrivlyData), ctx, userDeviceID, deviceDefinitionID, vin, forceSetAll)
}

// PullVincarioValuation mocks base method.
func (m *MockDeviceDefinitionService) PullVincarioValuation(ctx context.Context, userDeiceID, deviceDefinitionID, vin string) (services.DrivlyDataStatusEnum, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullVincarioValuation", ctx, userDeiceID, deviceDefinitionID, vin)
	ret0, _ := ret[0].(services.DrivlyDataStatusEnum)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PullVincarioValuation indicates an expected call of PullVincarioValuation.
func (mr *MockDeviceDefinitionServiceMockRecorder) PullVincarioValuation(ctx, userDeiceID, deviceDefinitionID, vin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullVincarioValuation", reflect.TypeOf((*MockDeviceDefinitionService)(nil).PullVincarioValuation), ctx, userDeiceID, deviceDefinitionID, vin)
}

// UpdateDeviceDefinitionFromNHTSA mocks base method.
func (m *MockDeviceDefinitionService) UpdateDeviceDefinitionFromNHTSA(ctx context.Context, deviceDefinitionID, vin string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDeviceDefinitionFromNHTSA", ctx, deviceDefinitionID, vin)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDeviceDefinitionFromNHTSA indicates an expected call of UpdateDeviceDefinitionFromNHTSA.
func (mr *MockDeviceDefinitionServiceMockRecorder) UpdateDeviceDefinitionFromNHTSA(ctx, deviceDefinitionID, vin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDeviceDefinitionFromNHTSA", reflect.TypeOf((*MockDeviceDefinitionService)(nil).UpdateDeviceDefinitionFromNHTSA), ctx, deviceDefinitionID, vin)
}
