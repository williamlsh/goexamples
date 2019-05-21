// Code generated by MockGen. DO NOT EDIT.
// Source: goexamples/routeguide (interfaces: RouteGuideClient,RouteGuide_RouteChatClient)

// Package mock_routeguide is a generated GoMock package.
package mock_routeguide

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	routeguide "goexamples/routeguide"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	reflect "reflect"
)

// MockRouteGuideClient is a mock of RouteGuideClient interface
type MockRouteGuideClient struct {
	ctrl     *gomock.Controller
	recorder *MockRouteGuideClientMockRecorder
}

// MockRouteGuideClientMockRecorder is the mock recorder for MockRouteGuideClient
type MockRouteGuideClientMockRecorder struct {
	mock *MockRouteGuideClient
}

// NewMockRouteGuideClient creates a new mock instance
func NewMockRouteGuideClient(ctrl *gomock.Controller) *MockRouteGuideClient {
	mock := &MockRouteGuideClient{ctrl: ctrl}
	mock.recorder = &MockRouteGuideClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRouteGuideClient) EXPECT() *MockRouteGuideClientMockRecorder {
	return m.recorder
}

// GetFeature mocks base method
func (m *MockRouteGuideClient) GetFeature(arg0 context.Context, arg1 *routeguide.Point, arg2 ...grpc.CallOption) (*routeguide.Feature, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFeature", varargs...)
	ret0, _ := ret[0].(*routeguide.Feature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeature indicates an expected call of GetFeature
func (mr *MockRouteGuideClientMockRecorder) GetFeature(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeature", reflect.TypeOf((*MockRouteGuideClient)(nil).GetFeature), varargs...)
}

// ListFeatures mocks base method
func (m *MockRouteGuideClient) ListFeatures(arg0 context.Context, arg1 *routeguide.Rectangle, arg2 ...grpc.CallOption) (routeguide.RouteGuide_ListFeaturesClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFeatures", varargs...)
	ret0, _ := ret[0].(routeguide.RouteGuide_ListFeaturesClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeatures indicates an expected call of ListFeatures
func (mr *MockRouteGuideClientMockRecorder) ListFeatures(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeatures", reflect.TypeOf((*MockRouteGuideClient)(nil).ListFeatures), varargs...)
}

// RecordRoute mocks base method
func (m *MockRouteGuideClient) RecordRoute(arg0 context.Context, arg1 ...grpc.CallOption) (routeguide.RouteGuide_RecordRouteClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RecordRoute", varargs...)
	ret0, _ := ret[0].(routeguide.RouteGuide_RecordRouteClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RecordRoute indicates an expected call of RecordRoute
func (mr *MockRouteGuideClientMockRecorder) RecordRoute(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordRoute", reflect.TypeOf((*MockRouteGuideClient)(nil).RecordRoute), varargs...)
}

// RouteChat mocks base method
func (m *MockRouteGuideClient) RouteChat(arg0 context.Context, arg1 ...grpc.CallOption) (routeguide.RouteGuide_RouteChatClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RouteChat", varargs...)
	ret0, _ := ret[0].(routeguide.RouteGuide_RouteChatClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RouteChat indicates an expected call of RouteChat
func (mr *MockRouteGuideClientMockRecorder) RouteChat(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RouteChat", reflect.TypeOf((*MockRouteGuideClient)(nil).RouteChat), varargs...)
}

// MockRouteGuide_RouteChatClient is a mock of RouteGuide_RouteChatClient interface
type MockRouteGuide_RouteChatClient struct {
	ctrl     *gomock.Controller
	recorder *MockRouteGuide_RouteChatClientMockRecorder
}

// MockRouteGuide_RouteChatClientMockRecorder is the mock recorder for MockRouteGuide_RouteChatClient
type MockRouteGuide_RouteChatClientMockRecorder struct {
	mock *MockRouteGuide_RouteChatClient
}

// NewMockRouteGuide_RouteChatClient creates a new mock instance
func NewMockRouteGuide_RouteChatClient(ctrl *gomock.Controller) *MockRouteGuide_RouteChatClient {
	mock := &MockRouteGuide_RouteChatClient{ctrl: ctrl}
	mock.recorder = &MockRouteGuide_RouteChatClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRouteGuide_RouteChatClient) EXPECT() *MockRouteGuide_RouteChatClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockRouteGuide_RouteChatClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockRouteGuide_RouteChatClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockRouteGuide_RouteChatClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockRouteGuide_RouteChatClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).Context))
}

// Header mocks base method
func (m *MockRouteGuide_RouteChatClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockRouteGuide_RouteChatClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).Header))
}

// Recv mocks base method
func (m *MockRouteGuide_RouteChatClient) Recv() (*routeguide.RouteNote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*routeguide.RouteNote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockRouteGuide_RouteChatClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockRouteGuide_RouteChatClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockRouteGuide_RouteChatClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockRouteGuide_RouteChatClient) Send(arg0 *routeguide.RouteNote) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockRouteGuide_RouteChatClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).Send), arg0)
}

// SendMsg mocks base method
func (m *MockRouteGuide_RouteChatClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockRouteGuide_RouteChatClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockRouteGuide_RouteChatClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockRouteGuide_RouteChatClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockRouteGuide_RouteChatClient)(nil).Trailer))
}
