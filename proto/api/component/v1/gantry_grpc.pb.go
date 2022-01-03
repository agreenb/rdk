// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GantryServiceClient is the client API for GantryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GantryServiceClient interface {
	// CurrentPosition gets the current position of a gantry of the underlying robot.
	CurrentPosition(ctx context.Context, in *GantryServiceCurrentPositionRequest, opts ...grpc.CallOption) (*GantryServiceCurrentPositionResponse, error)
	// MoveToPosition moves a gantry of the underlying robot to the requested position.
	MoveToPosition(ctx context.Context, in *GantryServiceMoveToPositionRequest, opts ...grpc.CallOption) (*GantryServiceMoveToPositionResponse, error)
	// Lengths gets the lengths of a gantry of the underlying robot.
	Lengths(ctx context.Context, in *GantryServiceLengthsRequest, opts ...grpc.CallOption) (*GantryServiceLengthsResponse, error)
}

type gantryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGantryServiceClient(cc grpc.ClientConnInterface) GantryServiceClient {
	return &gantryServiceClient{cc}
}

func (c *gantryServiceClient) CurrentPosition(ctx context.Context, in *GantryServiceCurrentPositionRequest, opts ...grpc.CallOption) (*GantryServiceCurrentPositionResponse, error) {
	out := new(GantryServiceCurrentPositionResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.GantryService/CurrentPosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gantryServiceClient) MoveToPosition(ctx context.Context, in *GantryServiceMoveToPositionRequest, opts ...grpc.CallOption) (*GantryServiceMoveToPositionResponse, error) {
	out := new(GantryServiceMoveToPositionResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.GantryService/MoveToPosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gantryServiceClient) Lengths(ctx context.Context, in *GantryServiceLengthsRequest, opts ...grpc.CallOption) (*GantryServiceLengthsResponse, error) {
	out := new(GantryServiceLengthsResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.GantryService/Lengths", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GantryServiceServer is the server API for GantryService service.
// All implementations must embed UnimplementedGantryServiceServer
// for forward compatibility
type GantryServiceServer interface {
	// CurrentPosition gets the current position of a gantry of the underlying robot.
	CurrentPosition(context.Context, *GantryServiceCurrentPositionRequest) (*GantryServiceCurrentPositionResponse, error)
	// MoveToPosition moves a gantry of the underlying robot to the requested position.
	MoveToPosition(context.Context, *GantryServiceMoveToPositionRequest) (*GantryServiceMoveToPositionResponse, error)
	// Lengths gets the lengths of a gantry of the underlying robot.
	Lengths(context.Context, *GantryServiceLengthsRequest) (*GantryServiceLengthsResponse, error)
	mustEmbedUnimplementedGantryServiceServer()
}

// UnimplementedGantryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGantryServiceServer struct {
}

func (UnimplementedGantryServiceServer) CurrentPosition(context.Context, *GantryServiceCurrentPositionRequest) (*GantryServiceCurrentPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CurrentPosition not implemented")
}
func (UnimplementedGantryServiceServer) MoveToPosition(context.Context, *GantryServiceMoveToPositionRequest) (*GantryServiceMoveToPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveToPosition not implemented")
}
func (UnimplementedGantryServiceServer) Lengths(context.Context, *GantryServiceLengthsRequest) (*GantryServiceLengthsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lengths not implemented")
}
func (UnimplementedGantryServiceServer) mustEmbedUnimplementedGantryServiceServer() {}

// UnsafeGantryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GantryServiceServer will
// result in compilation errors.
type UnsafeGantryServiceServer interface {
	mustEmbedUnimplementedGantryServiceServer()
}

func RegisterGantryServiceServer(s grpc.ServiceRegistrar, srv GantryServiceServer) {
	s.RegisterService(&GantryService_ServiceDesc, srv)
}

func _GantryService_CurrentPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GantryServiceCurrentPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GantryServiceServer).CurrentPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.GantryService/CurrentPosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GantryServiceServer).CurrentPosition(ctx, req.(*GantryServiceCurrentPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GantryService_MoveToPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GantryServiceMoveToPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GantryServiceServer).MoveToPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.GantryService/MoveToPosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GantryServiceServer).MoveToPosition(ctx, req.(*GantryServiceMoveToPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GantryService_Lengths_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GantryServiceLengthsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GantryServiceServer).Lengths(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.GantryService/Lengths",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GantryServiceServer).Lengths(ctx, req.(*GantryServiceLengthsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GantryService_ServiceDesc is the grpc.ServiceDesc for GantryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GantryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.api.component.v1.GantryService",
	HandlerType: (*GantryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CurrentPosition",
			Handler:    _GantryService_CurrentPosition_Handler,
		},
		{
			MethodName: "MoveToPosition",
			Handler:    _GantryService_MoveToPosition_Handler,
		},
		{
			MethodName: "Lengths",
			Handler:    _GantryService_Lengths_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/api/component/v1/gantry.proto",
}
