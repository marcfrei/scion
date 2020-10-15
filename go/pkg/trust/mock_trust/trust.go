// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/pkg/trust (interfaces: DB,Fetcher,Inspector,KeyRing,Provider,Recurser,Router,X509KeyPairLoader)

// Package mock_trust is a generated GoMock package.
package mock_trust

import (
	context "context"
	crypto "crypto"
	tls "crypto/tls"
	x509 "crypto/x509"
	gomock "github.com/golang/mock/gomock"
	addr "github.com/scionproto/scion/go/lib/addr"
	cppki "github.com/scionproto/scion/go/lib/scrypto/cppki"
	trust "github.com/scionproto/scion/go/pkg/trust"
	net "net"
	reflect "reflect"
)

// MockDB is a mock of DB interface
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// Chains mocks base method
func (m *MockDB) Chains(arg0 context.Context, arg1 trust.ChainQuery) ([][]*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chains", arg0, arg1)
	ret0, _ := ret[0].([][]*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Chains indicates an expected call of Chains
func (mr *MockDBMockRecorder) Chains(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chains", reflect.TypeOf((*MockDB)(nil).Chains), arg0, arg1)
}

// Close mocks base method
func (m *MockDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDB)(nil).Close))
}

// InsertChain mocks base method
func (m *MockDB) InsertChain(arg0 context.Context, arg1 []*x509.Certificate) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertChain", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertChain indicates an expected call of InsertChain
func (mr *MockDBMockRecorder) InsertChain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertChain", reflect.TypeOf((*MockDB)(nil).InsertChain), arg0, arg1)
}

// InsertTRC mocks base method
func (m *MockDB) InsertTRC(arg0 context.Context, arg1 cppki.SignedTRC) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTRC", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTRC indicates an expected call of InsertTRC
func (mr *MockDBMockRecorder) InsertTRC(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTRC", reflect.TypeOf((*MockDB)(nil).InsertTRC), arg0, arg1)
}

// SetMaxIdleConns mocks base method
func (m *MockDB) SetMaxIdleConns(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxIdleConns", arg0)
}

// SetMaxIdleConns indicates an expected call of SetMaxIdleConns
func (mr *MockDBMockRecorder) SetMaxIdleConns(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxIdleConns", reflect.TypeOf((*MockDB)(nil).SetMaxIdleConns), arg0)
}

// SetMaxOpenConns mocks base method
func (m *MockDB) SetMaxOpenConns(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxOpenConns", arg0)
}

// SetMaxOpenConns indicates an expected call of SetMaxOpenConns
func (mr *MockDBMockRecorder) SetMaxOpenConns(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxOpenConns", reflect.TypeOf((*MockDB)(nil).SetMaxOpenConns), arg0)
}

// SignedTRC mocks base method
func (m *MockDB) SignedTRC(arg0 context.Context, arg1 cppki.TRCID) (cppki.SignedTRC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignedTRC", arg0, arg1)
	ret0, _ := ret[0].(cppki.SignedTRC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignedTRC indicates an expected call of SignedTRC
func (mr *MockDBMockRecorder) SignedTRC(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignedTRC", reflect.TypeOf((*MockDB)(nil).SignedTRC), arg0, arg1)
}

// MockFetcher is a mock of Fetcher interface
type MockFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockFetcherMockRecorder
}

// MockFetcherMockRecorder is the mock recorder for MockFetcher
type MockFetcherMockRecorder struct {
	mock *MockFetcher
}

// NewMockFetcher creates a new mock instance
func NewMockFetcher(ctrl *gomock.Controller) *MockFetcher {
	mock := &MockFetcher{ctrl: ctrl}
	mock.recorder = &MockFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFetcher) EXPECT() *MockFetcherMockRecorder {
	return m.recorder
}

// Chains mocks base method
func (m *MockFetcher) Chains(arg0 context.Context, arg1 trust.ChainQuery, arg2 net.Addr) ([][]*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chains", arg0, arg1, arg2)
	ret0, _ := ret[0].([][]*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Chains indicates an expected call of Chains
func (mr *MockFetcherMockRecorder) Chains(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chains", reflect.TypeOf((*MockFetcher)(nil).Chains), arg0, arg1, arg2)
}

// TRC mocks base method
func (m *MockFetcher) TRC(arg0 context.Context, arg1 cppki.TRCID, arg2 net.Addr) (cppki.SignedTRC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TRC", arg0, arg1, arg2)
	ret0, _ := ret[0].(cppki.SignedTRC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TRC indicates an expected call of TRC
func (mr *MockFetcherMockRecorder) TRC(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TRC", reflect.TypeOf((*MockFetcher)(nil).TRC), arg0, arg1, arg2)
}

// MockInspector is a mock of Inspector interface
type MockInspector struct {
	ctrl     *gomock.Controller
	recorder *MockInspectorMockRecorder
}

// MockInspectorMockRecorder is the mock recorder for MockInspector
type MockInspectorMockRecorder struct {
	mock *MockInspector
}

// NewMockInspector creates a new mock instance
func NewMockInspector(ctrl *gomock.Controller) *MockInspector {
	mock := &MockInspector{ctrl: ctrl}
	mock.recorder = &MockInspectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInspector) EXPECT() *MockInspectorMockRecorder {
	return m.recorder
}

// ByAttributes mocks base method
func (m *MockInspector) ByAttributes(arg0 context.Context, arg1 addr.ISD, arg2 trust.Attribute) ([]addr.IA, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByAttributes", arg0, arg1, arg2)
	ret0, _ := ret[0].([]addr.IA)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByAttributes indicates an expected call of ByAttributes
func (mr *MockInspectorMockRecorder) ByAttributes(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByAttributes", reflect.TypeOf((*MockInspector)(nil).ByAttributes), arg0, arg1, arg2)
}

// HasAttributes mocks base method
func (m *MockInspector) HasAttributes(arg0 context.Context, arg1 addr.IA, arg2 trust.Attribute) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasAttributes", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasAttributes indicates an expected call of HasAttributes
func (mr *MockInspectorMockRecorder) HasAttributes(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasAttributes", reflect.TypeOf((*MockInspector)(nil).HasAttributes), arg0, arg1, arg2)
}

// MockKeyRing is a mock of KeyRing interface
type MockKeyRing struct {
	ctrl     *gomock.Controller
	recorder *MockKeyRingMockRecorder
}

// MockKeyRingMockRecorder is the mock recorder for MockKeyRing
type MockKeyRingMockRecorder struct {
	mock *MockKeyRing
}

// NewMockKeyRing creates a new mock instance
func NewMockKeyRing(ctrl *gomock.Controller) *MockKeyRing {
	mock := &MockKeyRing{ctrl: ctrl}
	mock.recorder = &MockKeyRingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyRing) EXPECT() *MockKeyRingMockRecorder {
	return m.recorder
}

// PrivateKeys mocks base method
func (m *MockKeyRing) PrivateKeys(arg0 context.Context) ([]crypto.Signer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrivateKeys", arg0)
	ret0, _ := ret[0].([]crypto.Signer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrivateKeys indicates an expected call of PrivateKeys
func (mr *MockKeyRingMockRecorder) PrivateKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrivateKeys", reflect.TypeOf((*MockKeyRing)(nil).PrivateKeys), arg0)
}

// MockProvider is a mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// GetChains mocks base method
func (m *MockProvider) GetChains(arg0 context.Context, arg1 trust.ChainQuery, arg2 ...trust.Option) ([][]*x509.Certificate, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetChains", varargs...)
	ret0, _ := ret[0].([][]*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChains indicates an expected call of GetChains
func (mr *MockProviderMockRecorder) GetChains(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChains", reflect.TypeOf((*MockProvider)(nil).GetChains), varargs...)
}

// GetSignedTRC mocks base method
func (m *MockProvider) GetSignedTRC(arg0 context.Context, arg1 cppki.TRCID, arg2 ...trust.Option) (cppki.SignedTRC, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSignedTRC", varargs...)
	ret0, _ := ret[0].(cppki.SignedTRC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSignedTRC indicates an expected call of GetSignedTRC
func (mr *MockProviderMockRecorder) GetSignedTRC(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSignedTRC", reflect.TypeOf((*MockProvider)(nil).GetSignedTRC), varargs...)
}

// NotifyTRC mocks base method
func (m *MockProvider) NotifyTRC(arg0 context.Context, arg1 cppki.TRCID, arg2 ...trust.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NotifyTRC", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyTRC indicates an expected call of NotifyTRC
func (mr *MockProviderMockRecorder) NotifyTRC(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyTRC", reflect.TypeOf((*MockProvider)(nil).NotifyTRC), varargs...)
}

// MockRecurser is a mock of Recurser interface
type MockRecurser struct {
	ctrl     *gomock.Controller
	recorder *MockRecurserMockRecorder
}

// MockRecurserMockRecorder is the mock recorder for MockRecurser
type MockRecurserMockRecorder struct {
	mock *MockRecurser
}

// NewMockRecurser creates a new mock instance
func NewMockRecurser(ctrl *gomock.Controller) *MockRecurser {
	mock := &MockRecurser{ctrl: ctrl}
	mock.recorder = &MockRecurserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRecurser) EXPECT() *MockRecurserMockRecorder {
	return m.recorder
}

// AllowRecursion mocks base method
func (m *MockRecurser) AllowRecursion(arg0 net.Addr) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowRecursion", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AllowRecursion indicates an expected call of AllowRecursion
func (mr *MockRecurserMockRecorder) AllowRecursion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowRecursion", reflect.TypeOf((*MockRecurser)(nil).AllowRecursion), arg0)
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

// ChooseServer mocks base method
func (m *MockRouter) ChooseServer(arg0 context.Context, arg1 addr.ISD) (net.Addr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChooseServer", arg0, arg1)
	ret0, _ := ret[0].(net.Addr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChooseServer indicates an expected call of ChooseServer
func (mr *MockRouterMockRecorder) ChooseServer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChooseServer", reflect.TypeOf((*MockRouter)(nil).ChooseServer), arg0, arg1)
}

// MockX509KeyPairLoader is a mock of X509KeyPairLoader interface
type MockX509KeyPairLoader struct {
	ctrl     *gomock.Controller
	recorder *MockX509KeyPairLoaderMockRecorder
}

// MockX509KeyPairLoaderMockRecorder is the mock recorder for MockX509KeyPairLoader
type MockX509KeyPairLoaderMockRecorder struct {
	mock *MockX509KeyPairLoader
}

// NewMockX509KeyPairLoader creates a new mock instance
func NewMockX509KeyPairLoader(ctrl *gomock.Controller) *MockX509KeyPairLoader {
	mock := &MockX509KeyPairLoader{ctrl: ctrl}
	mock.recorder = &MockX509KeyPairLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockX509KeyPairLoader) EXPECT() *MockX509KeyPairLoaderMockRecorder {
	return m.recorder
}

// LoadX509KeyPair mocks base method
func (m *MockX509KeyPairLoader) LoadX509KeyPair() (*tls.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadX509KeyPair")
	ret0, _ := ret[0].(*tls.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadX509KeyPair indicates an expected call of LoadX509KeyPair
func (mr *MockX509KeyPairLoaderMockRecorder) LoadX509KeyPair() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadX509KeyPair", reflect.TypeOf((*MockX509KeyPairLoader)(nil).LoadX509KeyPair))
}
