// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azure/radius/pkg/radrp/backend/deployment (interfaces: DeploymentProcessor)

// Package deployment is a generated GoMock package.
package deployment

import (
	context "context"
	reflect "reflect"

	azresources "github.com/Azure/radius/pkg/azure/azresources"
	db "github.com/Azure/radius/pkg/radrp/db"
	gomock "github.com/golang/mock/gomock"
)

// MockDeploymentProcessor is a mock of DeploymentProcessor interface.
type MockDeploymentProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockDeploymentProcessorMockRecorder
}

// MockDeploymentProcessorMockRecorder is the mock recorder for MockDeploymentProcessor.
type MockDeploymentProcessorMockRecorder struct {
	mock *MockDeploymentProcessor
}

// NewMockDeploymentProcessor creates a new mock instance.
func NewMockDeploymentProcessor(ctrl *gomock.Controller) *MockDeploymentProcessor {
	mock := &MockDeploymentProcessor{ctrl: ctrl}
	mock.recorder = &MockDeploymentProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeploymentProcessor) EXPECT() *MockDeploymentProcessorMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDeploymentProcessor) Delete(arg0 context.Context, arg1 azresources.ResourceID, arg2 db.RadiusResource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeploymentProcessorMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeploymentProcessor)(nil).Delete), arg0, arg1, arg2)
}

// Deploy mocks base method.
func (m *MockDeploymentProcessor) Deploy(arg0 context.Context, arg1 azresources.ResourceID, arg2 db.RadiusResource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deploy indicates an expected call of Deploy.
func (mr *MockDeploymentProcessorMockRecorder) Deploy(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockDeploymentProcessor)(nil).Deploy), arg0, arg1, arg2)
}

// FetchSecrets mocks base method.
func (m *MockDeploymentProcessor) FetchSecrets(arg0 context.Context, arg1 azresources.ResourceID, arg2 db.RadiusResource) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchSecrets", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchSecrets indicates an expected call of FetchSecrets.
func (mr *MockDeploymentProcessorMockRecorder) FetchSecrets(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchSecrets", reflect.TypeOf((*MockDeploymentProcessor)(nil).FetchSecrets), arg0, arg1, arg2)
}
