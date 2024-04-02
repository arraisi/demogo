// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/proto/product/product_grpc.pb.go

// Package mock_product is a generated GoMock package.
package mock_product

import (
	context "context"
	product "demogo/internal/proto/product"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockProductServiceClient is a mock of ProductServiceClient interface.
type MockProductServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockProductServiceClientMockRecorder
}

// MockProductServiceClientMockRecorder is the mock recorder for MockProductServiceClient.
type MockProductServiceClientMockRecorder struct {
	mock *MockProductServiceClient
}

// NewMockProductServiceClient creates a new mock instance.
func NewMockProductServiceClient(ctrl *gomock.Controller) *MockProductServiceClient {
	mock := &MockProductServiceClient{ctrl: ctrl}
	mock.recorder = &MockProductServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductServiceClient) EXPECT() *MockProductServiceClientMockRecorder {
	return m.recorder
}

// GetProductByID mocks base method.
func (m *MockProductServiceClient) GetProductByID(ctx context.Context, in *product.GetProductIDRequest, opts ...grpc.CallOption) (*product.GetProductIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProductByID", varargs...)
	ret0, _ := ret[0].(*product.GetProductIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByID indicates an expected call of GetProductByID.
func (mr *MockProductServiceClientMockRecorder) GetProductByID(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByID", reflect.TypeOf((*MockProductServiceClient)(nil).GetProductByID), varargs...)
}

// SaveProduct mocks base method.
func (m *MockProductServiceClient) SaveProduct(ctx context.Context, in *product.SaveProductRequest, opts ...grpc.CallOption) (*product.SaveProductResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SaveProduct", varargs...)
	ret0, _ := ret[0].(*product.SaveProductResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveProduct indicates an expected call of SaveProduct.
func (mr *MockProductServiceClientMockRecorder) SaveProduct(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveProduct", reflect.TypeOf((*MockProductServiceClient)(nil).SaveProduct), varargs...)
}

// MockProductServiceServer is a mock of ProductServiceServer interface.
type MockProductServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockProductServiceServerMockRecorder
}

// MockProductServiceServerMockRecorder is the mock recorder for MockProductServiceServer.
type MockProductServiceServerMockRecorder struct {
	mock *MockProductServiceServer
}

// NewMockProductServiceServer creates a new mock instance.
func NewMockProductServiceServer(ctrl *gomock.Controller) *MockProductServiceServer {
	mock := &MockProductServiceServer{ctrl: ctrl}
	mock.recorder = &MockProductServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductServiceServer) EXPECT() *MockProductServiceServerMockRecorder {
	return m.recorder
}

// GetProductByID mocks base method.
func (m *MockProductServiceServer) GetProductByID(arg0 context.Context, arg1 *product.GetProductIDRequest) (*product.GetProductIDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductByID", arg0, arg1)
	ret0, _ := ret[0].(*product.GetProductIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByID indicates an expected call of GetProductByID.
func (mr *MockProductServiceServerMockRecorder) GetProductByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByID", reflect.TypeOf((*MockProductServiceServer)(nil).GetProductByID), arg0, arg1)
}

// SaveProduct mocks base method.
func (m *MockProductServiceServer) SaveProduct(arg0 context.Context, arg1 *product.SaveProductRequest) (*product.SaveProductResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveProduct", arg0, arg1)
	ret0, _ := ret[0].(*product.SaveProductResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveProduct indicates an expected call of SaveProduct.
func (mr *MockProductServiceServerMockRecorder) SaveProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveProduct", reflect.TypeOf((*MockProductServiceServer)(nil).SaveProduct), arg0, arg1)
}

// mustEmbedUnimplementedProductServiceServer mocks base method.
func (m *MockProductServiceServer) mustEmbedUnimplementedProductServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedProductServiceServer")
}

// mustEmbedUnimplementedProductServiceServer indicates an expected call of mustEmbedUnimplementedProductServiceServer.
func (mr *MockProductServiceServerMockRecorder) mustEmbedUnimplementedProductServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedProductServiceServer", reflect.TypeOf((*MockProductServiceServer)(nil).mustEmbedUnimplementedProductServiceServer))
}

// MockUnsafeProductServiceServer is a mock of UnsafeProductServiceServer interface.
type MockUnsafeProductServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeProductServiceServerMockRecorder
}

// MockUnsafeProductServiceServerMockRecorder is the mock recorder for MockUnsafeProductServiceServer.
type MockUnsafeProductServiceServerMockRecorder struct {
	mock *MockUnsafeProductServiceServer
}

// NewMockUnsafeProductServiceServer creates a new mock instance.
func NewMockUnsafeProductServiceServer(ctrl *gomock.Controller) *MockUnsafeProductServiceServer {
	mock := &MockUnsafeProductServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeProductServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeProductServiceServer) EXPECT() *MockUnsafeProductServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedProductServiceServer mocks base method.
func (m *MockUnsafeProductServiceServer) mustEmbedUnimplementedProductServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedProductServiceServer")
}

// mustEmbedUnimplementedProductServiceServer indicates an expected call of mustEmbedUnimplementedProductServiceServer.
func (mr *MockUnsafeProductServiceServerMockRecorder) mustEmbedUnimplementedProductServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedProductServiceServer", reflect.TypeOf((*MockUnsafeProductServiceServer)(nil).mustEmbedUnimplementedProductServiceServer))
}
