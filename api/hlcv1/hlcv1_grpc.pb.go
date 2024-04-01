// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: hlcv1.proto

package hlcv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// HCLServiceClient is the client API for HCLService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HCLServiceClient interface {
	Get(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetResp, error)
	BatchGet(ctx context.Context, in *BatchGetReq, opts ...grpc.CallOption) (*BatchGetResp, error)
}

type hCLServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHCLServiceClient(cc grpc.ClientConnInterface) HCLServiceClient {
	return &hCLServiceClient{cc}
}

func (c *hCLServiceClient) Get(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetResp, error) {
	out := new(GetResp)
	err := c.cc.Invoke(ctx, "/hlcv1.HCLService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hCLServiceClient) BatchGet(ctx context.Context, in *BatchGetReq, opts ...grpc.CallOption) (*BatchGetResp, error) {
	out := new(BatchGetResp)
	err := c.cc.Invoke(ctx, "/hlcv1.HCLService/BatchGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HCLServiceServer is the server API for HCLService service.
// All implementations must embed UnimplementedHCLServiceServer
// for forward compatibility
type HCLServiceServer interface {
	Get(context.Context, *emptypb.Empty) (*GetResp, error)
	BatchGet(context.Context, *BatchGetReq) (*BatchGetResp, error)
	mustEmbedUnimplementedHCLServiceServer()
}

// UnimplementedHCLServiceServer must be embedded to have forward compatible implementations.
type UnimplementedHCLServiceServer struct {
}

func (UnimplementedHCLServiceServer) Get(context.Context, *emptypb.Empty) (*GetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedHCLServiceServer) BatchGet(context.Context, *BatchGetReq) (*BatchGetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchGet not implemented")
}
func (UnimplementedHCLServiceServer) mustEmbedUnimplementedHCLServiceServer() {}

// UnsafeHCLServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HCLServiceServer will
// result in compilation errors.
type UnsafeHCLServiceServer interface {
	mustEmbedUnimplementedHCLServiceServer()
}

func RegisterHCLServiceServer(s grpc.ServiceRegistrar, srv HCLServiceServer) {
	s.RegisterService(&HCLService_ServiceDesc, srv)
}

func _HCLService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HCLServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hlcv1.HCLService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HCLServiceServer).Get(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _HCLService_BatchGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchGetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HCLServiceServer).BatchGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hlcv1.HCLService/BatchGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HCLServiceServer).BatchGet(ctx, req.(*BatchGetReq))
	}
	return interceptor(ctx, in, info, handler)
}

// HCLService_ServiceDesc is the grpc.ServiceDesc for HCLService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HCLService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hlcv1.HCLService",
	HandlerType: (*HCLServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _HCLService_Get_Handler,
		},
		{
			MethodName: "BatchGet",
			Handler:    _HCLService_BatchGet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hlcv1.proto",
}