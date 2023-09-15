// Code generated by MockGen. DO NOT EDIT.
// Source: state.go

// Package operator is a generated GoMock package.
package operator

import (
	reflect "reflect"

	wire "github.com/bloxapp/ssv-dkg/pkgs/wire"
	gomock "github.com/golang/mock/gomock"
)

// MockInstance is a mock of Instance interface.
type MockInstance struct {
	ctrl     *gomock.Controller
	recorder *MockInstanceMockRecorder
}

// MockInstanceMockRecorder is the mock recorder for MockInstance.
type MockInstanceMockRecorder struct {
	mock *MockInstance
}

// NewMockInstance creates a new mock instance.
func NewMockInstance(ctrl *gomock.Controller) *MockInstance {
	mock := &MockInstance{ctrl: ctrl}
	mock.recorder = &MockInstanceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstance) EXPECT() *MockInstanceMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockInstance) Process(arg0 uint64, arg1 *wire.SignedTransport) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockInstanceMockRecorder) Process(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockInstance)(nil).Process), arg0, arg1)
}

// ReadError mocks base method.
func (m *MockInstance) ReadError() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadError")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadError indicates an expected call of ReadError.
func (mr *MockInstanceMockRecorder) ReadError() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadError", reflect.TypeOf((*MockInstance)(nil).ReadError))
}

// ReadResponse mocks base method.
func (m *MockInstance) ReadResponse() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadResponse")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// ReadResponse indicates an expected call of ReadResponse.
func (mr *MockInstanceMockRecorder) ReadResponse() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadResponse", reflect.TypeOf((*MockInstance)(nil).ReadResponse))
}

// VerifyInitiatorMessage mocks base method.
func (m *MockInstance) VerifyInitiatorMessage(msg, sig []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyInitiatorMessage", msg, sig)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyInitiatorMessage indicates an expected call of VerifyInitiatorMessage.
func (mr *MockInstanceMockRecorder) VerifyInitiatorMessage(msg, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyInitiatorMessage", reflect.TypeOf((*MockInstance)(nil).VerifyInitiatorMessage), msg, sig)
}
