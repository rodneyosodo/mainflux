// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: users/policies/auth.proto

package policies

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

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	Authorize(ctx context.Context, in *AuthorizeReq, opts ...grpc.CallOption) (*AuthorizeRes, error)
	Issue(ctx context.Context, in *IssueReq, opts ...grpc.CallOption) (*Token, error)
	Identify(ctx context.Context, in *Token, opts ...grpc.CallOption) (*UserIdentity, error)
	AddPolicy(ctx context.Context, in *AddPolicyReq, opts ...grpc.CallOption) (*AddPolicyRes, error)
	DeletePolicy(ctx context.Context, in *DeletePolicyReq, opts ...grpc.CallOption) (*DeletePolicyRes, error)
	ListPolicies(ctx context.Context, in *ListPoliciesReq, opts ...grpc.CallOption) (*ListPoliciesRes, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Authorize(ctx context.Context, in *AuthorizeReq, opts ...grpc.CallOption) (*AuthorizeRes, error) {
	out := new(AuthorizeRes)
	err := c.cc.Invoke(ctx, "/mainflux.users.policies.AuthService/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Issue(ctx context.Context, in *IssueReq, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/mainflux.users.policies.AuthService/Issue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Identify(ctx context.Context, in *Token, opts ...grpc.CallOption) (*UserIdentity, error) {
	out := new(UserIdentity)
	err := c.cc.Invoke(ctx, "/mainflux.users.policies.AuthService/Identify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) AddPolicy(ctx context.Context, in *AddPolicyReq, opts ...grpc.CallOption) (*AddPolicyRes, error) {
	out := new(AddPolicyRes)
	err := c.cc.Invoke(ctx, "/mainflux.users.policies.AuthService/AddPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) DeletePolicy(ctx context.Context, in *DeletePolicyReq, opts ...grpc.CallOption) (*DeletePolicyRes, error) {
	out := new(DeletePolicyRes)
	err := c.cc.Invoke(ctx, "/mainflux.users.policies.AuthService/DeletePolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ListPolicies(ctx context.Context, in *ListPoliciesReq, opts ...grpc.CallOption) (*ListPoliciesRes, error) {
	out := new(ListPoliciesRes)
	err := c.cc.Invoke(ctx, "/mainflux.users.policies.AuthService/ListPolicies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	Authorize(context.Context, *AuthorizeReq) (*AuthorizeRes, error)
	Issue(context.Context, *IssueReq) (*Token, error)
	Identify(context.Context, *Token) (*UserIdentity, error)
	AddPolicy(context.Context, *AddPolicyReq) (*AddPolicyRes, error)
	DeletePolicy(context.Context, *DeletePolicyReq) (*DeletePolicyRes, error)
	ListPolicies(context.Context, *ListPoliciesReq) (*ListPoliciesRes, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) Authorize(context.Context, *AuthorizeReq) (*AuthorizeRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (UnimplementedAuthServiceServer) Issue(context.Context, *IssueReq) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Issue not implemented")
}
func (UnimplementedAuthServiceServer) Identify(context.Context, *Token) (*UserIdentity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Identify not implemented")
}
func (UnimplementedAuthServiceServer) AddPolicy(context.Context, *AddPolicyReq) (*AddPolicyRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPolicy not implemented")
}
func (UnimplementedAuthServiceServer) DeletePolicy(context.Context, *DeletePolicyReq) (*DeletePolicyRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePolicy not implemented")
}
func (UnimplementedAuthServiceServer) ListPolicies(context.Context, *ListPoliciesReq) (*ListPoliciesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPolicies not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mainflux.users.policies.AuthService/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Authorize(ctx, req.(*AuthorizeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Issue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Issue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mainflux.users.policies.AuthService/Issue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Issue(ctx, req.(*IssueReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Identify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Identify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mainflux.users.policies.AuthService/Identify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Identify(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_AddPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPolicyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AddPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mainflux.users.policies.AuthService/AddPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AddPolicy(ctx, req.(*AddPolicyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_DeletePolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePolicyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).DeletePolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mainflux.users.policies.AuthService/DeletePolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).DeletePolicy(ctx, req.(*DeletePolicyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ListPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPoliciesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ListPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mainflux.users.policies.AuthService/ListPolicies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ListPolicies(ctx, req.(*ListPoliciesReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mainflux.users.policies.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _AuthService_Authorize_Handler,
		},
		{
			MethodName: "Issue",
			Handler:    _AuthService_Issue_Handler,
		},
		{
			MethodName: "Identify",
			Handler:    _AuthService_Identify_Handler,
		},
		{
			MethodName: "AddPolicy",
			Handler:    _AuthService_AddPolicy_Handler,
		},
		{
			MethodName: "DeletePolicy",
			Handler:    _AuthService_DeletePolicy_Handler,
		},
		{
			MethodName: "ListPolicies",
			Handler:    _AuthService_ListPolicies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users/policies/auth.proto",
}
