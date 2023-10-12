// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/private/app/appnet (interfaces: SVCResolver,Resolver)

// Package mock_infraenv is a generated GoMock package.
package mock_infraenv

import (
	context "context"
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	addr "github.com/scionproto/scion/pkg/addr"
	snet "github.com/scionproto/scion/pkg/snet"
	svc "github.com/scionproto/scion/private/svc"
)

// MockSVCResolver is a mock of SVCResolver interface.
type MockSVCResolver struct {
	ctrl     *gomock.Controller
	recorder *MockSVCResolverMockRecorder
}

// MockSVCResolverMockRecorder is the mock recorder for MockSVCResolver.
type MockSVCResolverMockRecorder struct {
	mock *MockSVCResolver
}

// NewMockSVCResolver creates a new mock instance.
func NewMockSVCResolver(ctrl *gomock.Controller) *MockSVCResolver {
	mock := &MockSVCResolver{ctrl: ctrl}
	mock.recorder = &MockSVCResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSVCResolver) EXPECT() *MockSVCResolverMockRecorder {
	return m.recorder
}

// GetUnderlay mocks base method.
func (m *MockSVCResolver) GetUnderlay(arg0 addr.SVC) (*net.UDPAddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnderlay", arg0)
	ret0, _ := ret[0].(*net.UDPAddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnderlay indicates an expected call of GetUnderlay.
func (mr *MockSVCResolverMockRecorder) GetUnderlay(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnderlay", reflect.TypeOf((*MockSVCResolver)(nil).GetUnderlay), arg0)
}

// MockResolver is a mock of Resolver interface.
type MockResolver struct {
	ctrl     *gomock.Controller
	recorder *MockResolverMockRecorder
}

// MockResolverMockRecorder is the mock recorder for MockResolver.
type MockResolverMockRecorder struct {
	mock *MockResolver
}

// NewMockResolver creates a new mock instance.
func NewMockResolver(ctrl *gomock.Controller) *MockResolver {
	mock := &MockResolver{ctrl: ctrl}
	mock.recorder = &MockResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResolver) EXPECT() *MockResolverMockRecorder {
	return m.recorder
}

// LookupSVC mocks base method.
func (m *MockResolver) LookupSVC(arg0 context.Context, arg1 snet.Path, arg2 addr.SVC) (*svc.Reply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupSVC", arg0, arg1, arg2)
	ret0, _ := ret[0].(*svc.Reply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupSVC indicates an expected call of LookupSVC.
func (mr *MockResolverMockRecorder) LookupSVC(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupSVC", reflect.TypeOf((*MockResolver)(nil).LookupSVC), arg0, arg1, arg2)
}