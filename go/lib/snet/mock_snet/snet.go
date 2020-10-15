// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/lib/snet (interfaces: PacketDispatcherService,Network,PacketConn,Path,PathMetadata,PathQuerier,Router,RevocationHandler)

// Package mock_snet is a generated GoMock package.
package mock_snet

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	addr "github.com/scionproto/scion/go/lib/addr"
	common "github.com/scionproto/scion/go/lib/common"
	snet "github.com/scionproto/scion/go/lib/snet"
	spath "github.com/scionproto/scion/go/lib/spath"
	net "net"
	reflect "reflect"
	time "time"
)

// MockPacketDispatcherService is a mock of PacketDispatcherService interface
type MockPacketDispatcherService struct {
	ctrl     *gomock.Controller
	recorder *MockPacketDispatcherServiceMockRecorder
}

// MockPacketDispatcherServiceMockRecorder is the mock recorder for MockPacketDispatcherService
type MockPacketDispatcherServiceMockRecorder struct {
	mock *MockPacketDispatcherService
}

// NewMockPacketDispatcherService creates a new mock instance
func NewMockPacketDispatcherService(ctrl *gomock.Controller) *MockPacketDispatcherService {
	mock := &MockPacketDispatcherService{ctrl: ctrl}
	mock.recorder = &MockPacketDispatcherServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPacketDispatcherService) EXPECT() *MockPacketDispatcherServiceMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockPacketDispatcherService) Register(arg0 context.Context, arg1 addr.IA, arg2 *net.UDPAddr, arg3 addr.HostSVC) (snet.PacketConn, uint16, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(snet.PacketConn)
	ret1, _ := ret[1].(uint16)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Register indicates an expected call of Register
func (mr *MockPacketDispatcherServiceMockRecorder) Register(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockPacketDispatcherService)(nil).Register), arg0, arg1, arg2, arg3)
}

// MockNetwork is a mock of Network interface
type MockNetwork struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkMockRecorder
}

// MockNetworkMockRecorder is the mock recorder for MockNetwork
type MockNetworkMockRecorder struct {
	mock *MockNetwork
}

// NewMockNetwork creates a new mock instance
func NewMockNetwork(ctrl *gomock.Controller) *MockNetwork {
	mock := &MockNetwork{ctrl: ctrl}
	mock.recorder = &MockNetworkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetwork) EXPECT() *MockNetworkMockRecorder {
	return m.recorder
}

// Dial mocks base method
func (m *MockNetwork) Dial(arg0 context.Context, arg1 string, arg2 *net.UDPAddr, arg3 *snet.UDPAddr, arg4 addr.HostSVC) (*snet.Conn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dial", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*snet.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Dial indicates an expected call of Dial
func (mr *MockNetworkMockRecorder) Dial(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dial", reflect.TypeOf((*MockNetwork)(nil).Dial), arg0, arg1, arg2, arg3, arg4)
}

// Listen mocks base method
func (m *MockNetwork) Listen(arg0 context.Context, arg1 string, arg2 *net.UDPAddr, arg3 addr.HostSVC) (*snet.Conn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Listen", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*snet.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Listen indicates an expected call of Listen
func (mr *MockNetworkMockRecorder) Listen(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listen", reflect.TypeOf((*MockNetwork)(nil).Listen), arg0, arg1, arg2, arg3)
}

// MockPacketConn is a mock of PacketConn interface
type MockPacketConn struct {
	ctrl     *gomock.Controller
	recorder *MockPacketConnMockRecorder
}

// MockPacketConnMockRecorder is the mock recorder for MockPacketConn
type MockPacketConnMockRecorder struct {
	mock *MockPacketConn
}

// NewMockPacketConn creates a new mock instance
func NewMockPacketConn(ctrl *gomock.Controller) *MockPacketConn {
	mock := &MockPacketConn{ctrl: ctrl}
	mock.recorder = &MockPacketConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPacketConn) EXPECT() *MockPacketConnMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockPacketConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockPacketConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPacketConn)(nil).Close))
}

// ReadFrom mocks base method
func (m *MockPacketConn) ReadFrom(arg0 *snet.Packet, arg1 *net.UDPAddr) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFrom", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadFrom indicates an expected call of ReadFrom
func (mr *MockPacketConnMockRecorder) ReadFrom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFrom", reflect.TypeOf((*MockPacketConn)(nil).ReadFrom), arg0, arg1)
}

// SetDeadline mocks base method
func (m *MockPacketConn) SetDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDeadline indicates an expected call of SetDeadline
func (mr *MockPacketConnMockRecorder) SetDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeadline", reflect.TypeOf((*MockPacketConn)(nil).SetDeadline), arg0)
}

// SetReadDeadline mocks base method
func (m *MockPacketConn) SetReadDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline
func (mr *MockPacketConnMockRecorder) SetReadDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockPacketConn)(nil).SetReadDeadline), arg0)
}

// SetWriteDeadline mocks base method
func (m *MockPacketConn) SetWriteDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteDeadline indicates an expected call of SetWriteDeadline
func (mr *MockPacketConnMockRecorder) SetWriteDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteDeadline", reflect.TypeOf((*MockPacketConn)(nil).SetWriteDeadline), arg0)
}

// WriteTo mocks base method
func (m *MockPacketConn) WriteTo(arg0 *snet.Packet, arg1 *net.UDPAddr) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteTo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteTo indicates an expected call of WriteTo
func (mr *MockPacketConnMockRecorder) WriteTo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteTo", reflect.TypeOf((*MockPacketConn)(nil).WriteTo), arg0, arg1)
}

// MockPath is a mock of Path interface
type MockPath struct {
	ctrl     *gomock.Controller
	recorder *MockPathMockRecorder
}

// MockPathMockRecorder is the mock recorder for MockPath
type MockPathMockRecorder struct {
	mock *MockPath
}

// NewMockPath creates a new mock instance
func NewMockPath(ctrl *gomock.Controller) *MockPath {
	mock := &MockPath{ctrl: ctrl}
	mock.recorder = &MockPathMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPath) EXPECT() *MockPathMockRecorder {
	return m.recorder
}

// Copy mocks base method
func (m *MockPath) Copy() snet.Path {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy")
	ret0, _ := ret[0].(snet.Path)
	return ret0
}

// Copy indicates an expected call of Copy
func (mr *MockPathMockRecorder) Copy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockPath)(nil).Copy))
}

// Destination mocks base method
func (m *MockPath) Destination() addr.IA {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destination")
	ret0, _ := ret[0].(addr.IA)
	return ret0
}

// Destination indicates an expected call of Destination
func (mr *MockPathMockRecorder) Destination() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destination", reflect.TypeOf((*MockPath)(nil).Destination))
}

// Interfaces mocks base method
func (m *MockPath) Interfaces() []snet.PathInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Interfaces")
	ret0, _ := ret[0].([]snet.PathInterface)
	return ret0
}

// Interfaces indicates an expected call of Interfaces
func (mr *MockPathMockRecorder) Interfaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Interfaces", reflect.TypeOf((*MockPath)(nil).Interfaces))
}

// Metadata mocks base method
func (m *MockPath) Metadata() snet.PathMetadata {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Metadata")
	ret0, _ := ret[0].(snet.PathMetadata)
	return ret0
}

// Metadata indicates an expected call of Metadata
func (mr *MockPathMockRecorder) Metadata() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metadata", reflect.TypeOf((*MockPath)(nil).Metadata))
}

// Path mocks base method
func (m *MockPath) Path() *spath.Path {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path")
	ret0, _ := ret[0].(*spath.Path)
	return ret0
}

// Path indicates an expected call of Path
func (mr *MockPathMockRecorder) Path() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockPath)(nil).Path))
}

// UnderlayNextHop mocks base method
func (m *MockPath) UnderlayNextHop() *net.UDPAddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnderlayNextHop")
	ret0, _ := ret[0].(*net.UDPAddr)
	return ret0
}

// UnderlayNextHop indicates an expected call of UnderlayNextHop
func (mr *MockPathMockRecorder) UnderlayNextHop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnderlayNextHop", reflect.TypeOf((*MockPath)(nil).UnderlayNextHop))
}

// MockPathMetadata is a mock of PathMetadata interface
type MockPathMetadata struct {
	ctrl     *gomock.Controller
	recorder *MockPathMetadataMockRecorder
}

// MockPathMetadataMockRecorder is the mock recorder for MockPathMetadata
type MockPathMetadataMockRecorder struct {
	mock *MockPathMetadata
}

// NewMockPathMetadata creates a new mock instance
func NewMockPathMetadata(ctrl *gomock.Controller) *MockPathMetadata {
	mock := &MockPathMetadata{ctrl: ctrl}
	mock.recorder = &MockPathMetadataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPathMetadata) EXPECT() *MockPathMetadataMockRecorder {
	return m.recorder
}

// Expiry mocks base method
func (m *MockPathMetadata) Expiry() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expiry")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Expiry indicates an expected call of Expiry
func (mr *MockPathMetadataMockRecorder) Expiry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expiry", reflect.TypeOf((*MockPathMetadata)(nil).Expiry))
}

// MTU mocks base method
func (m *MockPathMetadata) MTU() uint16 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MTU")
	ret0, _ := ret[0].(uint16)
	return ret0
}

// MTU indicates an expected call of MTU
func (mr *MockPathMetadataMockRecorder) MTU() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MTU", reflect.TypeOf((*MockPathMetadata)(nil).MTU))
}

// MockPathQuerier is a mock of PathQuerier interface
type MockPathQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockPathQuerierMockRecorder
}

// MockPathQuerierMockRecorder is the mock recorder for MockPathQuerier
type MockPathQuerierMockRecorder struct {
	mock *MockPathQuerier
}

// NewMockPathQuerier creates a new mock instance
func NewMockPathQuerier(ctrl *gomock.Controller) *MockPathQuerier {
	mock := &MockPathQuerier{ctrl: ctrl}
	mock.recorder = &MockPathQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPathQuerier) EXPECT() *MockPathQuerierMockRecorder {
	return m.recorder
}

// Query mocks base method
func (m *MockPathQuerier) Query(arg0 context.Context, arg1 addr.IA) ([]snet.Path, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0, arg1)
	ret0, _ := ret[0].([]snet.Path)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockPathQuerierMockRecorder) Query(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockPathQuerier)(nil).Query), arg0, arg1)
}

// MockRouter is a mock of Router interface
type MockRouter struct {
	ctrl     *gomock.Controller
	recorder *MockRouterMockRecorder
}

// MockRouterMockRecorder is the mock recorder for MockRouter
type MockRouterMockRecorder struct {
	mock *MockRouter
}

// NewMockRouter creates a new mock instance
func NewMockRouter(ctrl *gomock.Controller) *MockRouter {
	mock := &MockRouter{ctrl: ctrl}
	mock.recorder = &MockRouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRouter) EXPECT() *MockRouterMockRecorder {
	return m.recorder
}

// AllRoutes mocks base method
func (m *MockRouter) AllRoutes(arg0 context.Context, arg1 addr.IA) ([]snet.Path, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllRoutes", arg0, arg1)
	ret0, _ := ret[0].([]snet.Path)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllRoutes indicates an expected call of AllRoutes
func (mr *MockRouterMockRecorder) AllRoutes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllRoutes", reflect.TypeOf((*MockRouter)(nil).AllRoutes), arg0, arg1)
}

// Route mocks base method
func (m *MockRouter) Route(arg0 context.Context, arg1 addr.IA) (snet.Path, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Route", arg0, arg1)
	ret0, _ := ret[0].(snet.Path)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Route indicates an expected call of Route
func (mr *MockRouterMockRecorder) Route(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Route", reflect.TypeOf((*MockRouter)(nil).Route), arg0, arg1)
}

// MockRevocationHandler is a mock of RevocationHandler interface
type MockRevocationHandler struct {
	ctrl     *gomock.Controller
	recorder *MockRevocationHandlerMockRecorder
}

// MockRevocationHandlerMockRecorder is the mock recorder for MockRevocationHandler
type MockRevocationHandlerMockRecorder struct {
	mock *MockRevocationHandler
}

// NewMockRevocationHandler creates a new mock instance
func NewMockRevocationHandler(ctrl *gomock.Controller) *MockRevocationHandler {
	mock := &MockRevocationHandler{ctrl: ctrl}
	mock.recorder = &MockRevocationHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRevocationHandler) EXPECT() *MockRevocationHandlerMockRecorder {
	return m.recorder
}

// RevokeRaw mocks base method
func (m *MockRevocationHandler) RevokeRaw(arg0 context.Context, arg1 common.RawBytes) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RevokeRaw", arg0, arg1)
}

// RevokeRaw indicates an expected call of RevokeRaw
func (mr *MockRevocationHandlerMockRecorder) RevokeRaw(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeRaw", reflect.TypeOf((*MockRevocationHandler)(nil).RevokeRaw), arg0, arg1)
}
