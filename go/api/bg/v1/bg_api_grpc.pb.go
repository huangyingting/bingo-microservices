// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: v1/bg_api.proto

package v1

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

// GeoClient is the client API for Geo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GeoClient interface {
	// Sends ip address to geo request
	Location(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*LocationReply, error)
	Liveness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusReply, error)
	Readiness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusReply, error)
}

type geoClient struct {
	cc grpc.ClientConnInterface
}

func NewGeoClient(cc grpc.ClientConnInterface) GeoClient {
	return &geoClient{cc}
}

func (c *geoClient) Location(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*LocationReply, error) {
	out := new(LocationReply)
	err := c.cc.Invoke(ctx, "/api.bg.v1.Geo/Location", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geoClient) Liveness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := c.cc.Invoke(ctx, "/api.bg.v1.Geo/Liveness", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geoClient) Readiness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := c.cc.Invoke(ctx, "/api.bg.v1.Geo/Readiness", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GeoServer is the server API for Geo service.
// All implementations must embed UnimplementedGeoServer
// for forward compatibility
type GeoServer interface {
	// Sends ip address to geo request
	Location(context.Context, *LocationRequest) (*LocationReply, error)
	Liveness(context.Context, *emptypb.Empty) (*StatusReply, error)
	Readiness(context.Context, *emptypb.Empty) (*StatusReply, error)
	mustEmbedUnimplementedGeoServer()
}

// UnimplementedGeoServer must be embedded to have forward compatible implementations.
type UnimplementedGeoServer struct {
}

func (UnimplementedGeoServer) Location(context.Context, *LocationRequest) (*LocationReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Location not implemented")
}
func (UnimplementedGeoServer) Liveness(context.Context, *emptypb.Empty) (*StatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Liveness not implemented")
}
func (UnimplementedGeoServer) Readiness(context.Context, *emptypb.Empty) (*StatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Readiness not implemented")
}
func (UnimplementedGeoServer) mustEmbedUnimplementedGeoServer() {}

// UnsafeGeoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GeoServer will
// result in compilation errors.
type UnsafeGeoServer interface {
	mustEmbedUnimplementedGeoServer()
}

func RegisterGeoServer(s grpc.ServiceRegistrar, srv GeoServer) {
	s.RegisterService(&Geo_ServiceDesc, srv)
}

func _Geo_Location_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServer).Location(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.bg.v1.Geo/Location",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServer).Location(ctx, req.(*LocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Geo_Liveness_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServer).Liveness(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.bg.v1.Geo/Liveness",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServer).Liveness(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Geo_Readiness_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServer).Readiness(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.bg.v1.Geo/Readiness",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServer).Readiness(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Geo_ServiceDesc is the grpc.ServiceDesc for Geo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Geo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.bg.v1.Geo",
	HandlerType: (*GeoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Location",
			Handler:    _Geo_Location_Handler,
		},
		{
			MethodName: "Liveness",
			Handler:    _Geo_Liveness_Handler,
		},
		{
			MethodName: "Readiness",
			Handler:    _Geo_Readiness_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/bg_api.proto",
}
