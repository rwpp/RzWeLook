// Code generated by MockGen. DO NOT EDIT.
// Source: producer.go
//
// Generated by this command:
//
//	mockgen -source=producer.go -package=evtmocks -destination=mocks/producer.mock.go Producer
//
// Package evtmocks is a generated GoMock package.
package evtmocks

import (
	context "context"
	"github.com/rwpp/RzWeLook/pkg/migrator/events"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// ProduceInconsistentEvent mocks base method.
func (m *MockProducer) ProduceInconsistentEvent(ctx context.Context, event events.InconsistentEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProduceInconsistentEvent", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProduceInconsistentEvent indicates an expected call of ProduceInconsistentEvent.
func (mr *MockProducerMockRecorder) ProduceInconsistentEvent(ctx, event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProduceInconsistentEvent", reflect.TypeOf((*MockProducer)(nil).ProduceInconsistentEvent), ctx, event)
}
