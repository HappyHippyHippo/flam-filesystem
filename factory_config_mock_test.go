package filesystem

import (
	"reflect"

	"github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

// FactoryConfigMock is a mock of FactoryConfigMock interface.
type FactoryConfigMock struct {
	ctrl     *gomock.Controller
	recorder *FactoryConfigMockRecorder
}

// FactoryConfigMockRecorder is the mock recorder for FactoryConfig.
type FactoryConfigMockRecorder struct {
	mock *FactoryConfigMock
}

// NewFactoryConfigMock creates a new mock instance.
func NewFactoryConfigMock(ctrl *gomock.Controller) *FactoryConfigMock {
	mock := &FactoryConfigMock{ctrl: ctrl}
	mock.recorder = &FactoryConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *FactoryConfigMock) EXPECT() *FactoryConfigMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *FactoryConfigMock) Get(path string, def ...any) flam.Bag {
	m.ctrl.T.Helper()
	varargs := append([]any{path}, def...)
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(flam.Bag)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *FactoryConfigMockRecorder) Get(path any, def ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*FactoryConfigMock)(nil).Get), varargs...)
}
