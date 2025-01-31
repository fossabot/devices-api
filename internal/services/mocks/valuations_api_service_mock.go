// Code generated by MockGen. DO NOT EDIT.
// Source: valuations_api_service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	grpc "github.com/DIMO-Network/valuations-api/pkg/grpc"
	gomock "go.uber.org/mock/gomock"
)

// MockValuationsAPIService is a mock of ValuationsAPIService interface.
type MockValuationsAPIService struct {
	ctrl     *gomock.Controller
	recorder *MockValuationsAPIServiceMockRecorder
}

// MockValuationsAPIServiceMockRecorder is the mock recorder for MockValuationsAPIService.
type MockValuationsAPIServiceMockRecorder struct {
	mock *MockValuationsAPIService
}

// NewMockValuationsAPIService creates a new mock instance.
func NewMockValuationsAPIService(ctrl *gomock.Controller) *MockValuationsAPIService {
	mock := &MockValuationsAPIService{ctrl: ctrl}
	mock.recorder = &MockValuationsAPIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValuationsAPIService) EXPECT() *MockValuationsAPIServiceMockRecorder {
	return m.recorder
}

// GetUserDeviceOffers mocks base method.
func (m *MockValuationsAPIService) GetUserDeviceOffers(ctx context.Context, userDeviceID string) (*grpc.DeviceOffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDeviceOffers", ctx, userDeviceID)
	ret0, _ := ret[0].(*grpc.DeviceOffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDeviceOffers indicates an expected call of GetUserDeviceOffers.
func (mr *MockValuationsAPIServiceMockRecorder) GetUserDeviceOffers(ctx, userDeviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDeviceOffers", reflect.TypeOf((*MockValuationsAPIService)(nil).GetUserDeviceOffers), ctx, userDeviceID)
}

// GetUserDeviceValuations mocks base method.
func (m *MockValuationsAPIService) GetUserDeviceValuations(ctx context.Context, userDeviceID string) (*grpc.DeviceValuation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDeviceValuations", ctx, userDeviceID)
	ret0, _ := ret[0].(*grpc.DeviceValuation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDeviceValuations indicates an expected call of GetUserDeviceValuations.
func (mr *MockValuationsAPIServiceMockRecorder) GetUserDeviceValuations(ctx, userDeviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDeviceValuations", reflect.TypeOf((*MockValuationsAPIService)(nil).GetUserDeviceValuations), ctx, userDeviceID)
}
