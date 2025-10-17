package filesystem

import (
	"reflect"

	"github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

// DiskCreatorMock is a mock of FileSystemDiskCreator interface.
type DiskCreatorMock struct {
	ctrl     *gomock.Controller
	recorder *DiskCreatorMockRecorder
}

// DiskCreatorMockRecorder is the mock recorder for DiskCreator.
type DiskCreatorMockRecorder struct {
	mock *DiskCreatorMock
}

// NewDiskCreatorMock creates a new mock instance.
func NewDiskCreatorMock(ctrl *gomock.Controller) *DiskCreatorMock {
	mock := &DiskCreatorMock{ctrl: ctrl}
	mock.recorder = &DiskCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *DiskCreatorMock) EXPECT() *DiskCreatorMockRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *DiskCreatorMock) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *DiskCreatorMockRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*DiskCreatorMock)(nil).Accept), config)
}

// Create mocks base method.
func (m *DiskCreatorMock) Create(config flam.Bag) (Disk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(Disk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *DiskCreatorMockRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*DiskCreatorMock)(nil).Create), config)
}
