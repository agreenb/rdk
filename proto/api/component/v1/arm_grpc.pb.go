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

// ArmServiceClient is the client API for ArmService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArmServiceClient interface {
	// CurrentPosition gets the current position the end of the robot's arm expressed as X,Y,Z,ox,oy,oz,theta
	CurrentPosition(ctx context.Context, in *ArmServiceCurrentPositionRequest, opts ...grpc.CallOption) (*ArmServiceCurrentPositionResponse, error)
	// MoveToPosition moves the mount point of the robot's end effector to the requested position.
	MoveToPosition(ctx context.Context, in *ArmServiceMoveToPositionRequest, opts ...grpc.CallOption) (*ArmServiceMoveToPositionResponse, error)
	// CurrentJointPositions lists the joint positions (in degrees) of every joint on a robot
	CurrentJointPositions(ctx context.Context, in *ArmServiceCurrentJointPositionsRequest, opts ...grpc.CallOption) (*ArmServiceCurrentJointPositionsResponse, error)
	// MoveToJointPositions moves every joint on a robot's arm to specified angles which are expressed in degrees
	MoveToJointPositions(ctx context.Context, in *ArmServiceMoveToJointPositionsRequest, opts ...grpc.CallOption) (*ArmServiceMoveToJointPositionsResponse, error)
	// JointMoveDelta moves a specific joint of a robot by the the specified number of degrees
	JointMoveDelta(ctx context.Context, in *ArmServiceJointMoveDeltaRequest, opts ...grpc.CallOption) (*ArmServiceJointMoveDeltaResponse, error)
}

type armServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewArmServiceClient(cc grpc.ClientConnInterface) ArmServiceClient {
	return &armServiceClient{cc}
}

func (c *armServiceClient) CurrentPosition(ctx context.Context, in *ArmServiceCurrentPositionRequest, opts ...grpc.CallOption) (*ArmServiceCurrentPositionResponse, error) {
	out := new(ArmServiceCurrentPositionResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.ArmService/CurrentPosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *armServiceClient) MoveToPosition(ctx context.Context, in *ArmServiceMoveToPositionRequest, opts ...grpc.CallOption) (*ArmServiceMoveToPositionResponse, error) {
	out := new(ArmServiceMoveToPositionResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.ArmService/MoveToPosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *armServiceClient) CurrentJointPositions(ctx context.Context, in *ArmServiceCurrentJointPositionsRequest, opts ...grpc.CallOption) (*ArmServiceCurrentJointPositionsResponse, error) {
	out := new(ArmServiceCurrentJointPositionsResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.ArmService/CurrentJointPositions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *armServiceClient) MoveToJointPositions(ctx context.Context, in *ArmServiceMoveToJointPositionsRequest, opts ...grpc.CallOption) (*ArmServiceMoveToJointPositionsResponse, error) {
	out := new(ArmServiceMoveToJointPositionsResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.ArmService/MoveToJointPositions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *armServiceClient) JointMoveDelta(ctx context.Context, in *ArmServiceJointMoveDeltaRequest, opts ...grpc.CallOption) (*ArmServiceJointMoveDeltaResponse, error) {
	out := new(ArmServiceJointMoveDeltaResponse)
	err := c.cc.Invoke(ctx, "/proto.api.component.v1.ArmService/JointMoveDelta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArmServiceServer is the server API for ArmService service.
// All implementations must embed UnimplementedArmServiceServer
// for forward compatibility
type ArmServiceServer interface {
	// CurrentPosition gets the current position the end of the robot's arm expressed as X,Y,Z,ox,oy,oz,theta
	CurrentPosition(context.Context, *ArmServiceCurrentPositionRequest) (*ArmServiceCurrentPositionResponse, error)
	// MoveToPosition moves the mount point of the robot's end effector to the requested position.
	MoveToPosition(context.Context, *ArmServiceMoveToPositionRequest) (*ArmServiceMoveToPositionResponse, error)
	// CurrentJointPositions lists the joint positions (in degrees) of every joint on a robot
	CurrentJointPositions(context.Context, *ArmServiceCurrentJointPositionsRequest) (*ArmServiceCurrentJointPositionsResponse, error)
	// MoveToJointPositions moves every joint on a robot's arm to specified angles which are expressed in degrees
	MoveToJointPositions(context.Context, *ArmServiceMoveToJointPositionsRequest) (*ArmServiceMoveToJointPositionsResponse, error)
	// JointMoveDelta moves a specific joint of a robot by the the specified number of degrees
	JointMoveDelta(context.Context, *ArmServiceJointMoveDeltaRequest) (*ArmServiceJointMoveDeltaResponse, error)
	mustEmbedUnimplementedArmServiceServer()
}

// UnimplementedArmServiceServer must be embedded to have forward compatible implementations.
type UnimplementedArmServiceServer struct {
}

func (UnimplementedArmServiceServer) CurrentPosition(context.Context, *ArmServiceCurrentPositionRequest) (*ArmServiceCurrentPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CurrentPosition not implemented")
}
func (UnimplementedArmServiceServer) MoveToPosition(context.Context, *ArmServiceMoveToPositionRequest) (*ArmServiceMoveToPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveToPosition not implemented")
}
func (UnimplementedArmServiceServer) CurrentJointPositions(context.Context, *ArmServiceCurrentJointPositionsRequest) (*ArmServiceCurrentJointPositionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CurrentJointPositions not implemented")
}
func (UnimplementedArmServiceServer) MoveToJointPositions(context.Context, *ArmServiceMoveToJointPositionsRequest) (*ArmServiceMoveToJointPositionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveToJointPositions not implemented")
}
func (UnimplementedArmServiceServer) JointMoveDelta(context.Context, *ArmServiceJointMoveDeltaRequest) (*ArmServiceJointMoveDeltaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JointMoveDelta not implemented")
}
func (UnimplementedArmServiceServer) mustEmbedUnimplementedArmServiceServer() {}

// UnsafeArmServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArmServiceServer will
// result in compilation errors.
type UnsafeArmServiceServer interface {
	mustEmbedUnimplementedArmServiceServer()
}

func RegisterArmServiceServer(s grpc.ServiceRegistrar, srv ArmServiceServer) {
	s.RegisterService(&ArmService_ServiceDesc, srv)
}

func _ArmService_CurrentPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArmServiceCurrentPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArmServiceServer).CurrentPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.ArmService/CurrentPosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArmServiceServer).CurrentPosition(ctx, req.(*ArmServiceCurrentPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArmService_MoveToPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArmServiceMoveToPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArmServiceServer).MoveToPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.ArmService/MoveToPosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArmServiceServer).MoveToPosition(ctx, req.(*ArmServiceMoveToPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArmService_CurrentJointPositions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArmServiceCurrentJointPositionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArmServiceServer).CurrentJointPositions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.ArmService/CurrentJointPositions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArmServiceServer).CurrentJointPositions(ctx, req.(*ArmServiceCurrentJointPositionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArmService_MoveToJointPositions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArmServiceMoveToJointPositionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArmServiceServer).MoveToJointPositions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.ArmService/MoveToJointPositions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArmServiceServer).MoveToJointPositions(ctx, req.(*ArmServiceMoveToJointPositionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArmService_JointMoveDelta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArmServiceJointMoveDeltaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArmServiceServer).JointMoveDelta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.api.component.v1.ArmService/JointMoveDelta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArmServiceServer).JointMoveDelta(ctx, req.(*ArmServiceJointMoveDeltaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ArmService_ServiceDesc is the grpc.ServiceDesc for ArmService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ArmService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.api.component.v1.ArmService",
	HandlerType: (*ArmServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CurrentPosition",
			Handler:    _ArmService_CurrentPosition_Handler,
		},
		{
			MethodName: "MoveToPosition",
			Handler:    _ArmService_MoveToPosition_Handler,
		},
		{
			MethodName: "CurrentJointPositions",
			Handler:    _ArmService_CurrentJointPositions_Handler,
		},
		{
			MethodName: "MoveToJointPositions",
			Handler:    _ArmService_MoveToJointPositions_Handler,
		},
		{
			MethodName: "JointMoveDelta",
			Handler:    _ArmService_JointMoveDelta_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/api/component/v1/arm.proto",
}
