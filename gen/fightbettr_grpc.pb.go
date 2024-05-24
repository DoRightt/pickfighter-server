// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: fightbettr.proto

package gen

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

const (
	AuthService_Register_FullMethodName        = "/AuthService/Register"
	AuthService_RegisterConfirm_FullMethodName = "/AuthService/RegisterConfirm"
	AuthService_Login_FullMethodName           = "/AuthService/Login"
	AuthService_Logout_FullMethodName          = "/AuthService/Logout"
	AuthService_PasswordReset_FullMethodName   = "/AuthService/PasswordReset"
	AuthService_PasswordRecover_FullMethodName = "/AuthService/PasswordRecover"
	AuthService_Profile_FullMethodName         = "/AuthService/Profile"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	RegisterConfirm(ctx context.Context, in *RegisterConfirmRequest, opts ...grpc.CallOption) (*RegisterConfirmResponse, error)
	Login(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
	PasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*PasswordResetResponse, error)
	PasswordRecover(ctx context.Context, in *PasswordRecoveryRequest, opts ...grpc.CallOption) (*PasswordRecoveryResponse, error)
	Profile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, AuthService_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) RegisterConfirm(ctx context.Context, in *RegisterConfirmRequest, opts ...grpc.CallOption) (*RegisterConfirmResponse, error) {
	out := new(RegisterConfirmResponse)
	err := c.cc.Invoke(ctx, AuthService_RegisterConfirm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Login(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := c.cc.Invoke(ctx, AuthService_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, AuthService_Logout_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) PasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*PasswordResetResponse, error) {
	out := new(PasswordResetResponse)
	err := c.cc.Invoke(ctx, AuthService_PasswordReset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) PasswordRecover(ctx context.Context, in *PasswordRecoveryRequest, opts ...grpc.CallOption) (*PasswordRecoveryResponse, error) {
	out := new(PasswordRecoveryResponse)
	err := c.cc.Invoke(ctx, AuthService_PasswordRecover_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Profile(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error) {
	out := new(ProfileResponse)
	err := c.cc.Invoke(ctx, AuthService_Profile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	RegisterConfirm(context.Context, *RegisterConfirmRequest) (*RegisterConfirmResponse, error)
	Login(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	PasswordReset(context.Context, *PasswordResetRequest) (*PasswordResetResponse, error)
	PasswordRecover(context.Context, *PasswordRecoveryRequest) (*PasswordRecoveryResponse, error)
	Profile(context.Context, *ProfileRequest) (*ProfileResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthServiceServer) RegisterConfirm(context.Context, *RegisterConfirmRequest) (*RegisterConfirmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterConfirm not implemented")
}
func (UnimplementedAuthServiceServer) Login(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAuthServiceServer) PasswordReset(context.Context, *PasswordResetRequest) (*PasswordResetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PasswordReset not implemented")
}
func (UnimplementedAuthServiceServer) PasswordRecover(context.Context, *PasswordRecoveryRequest) (*PasswordRecoveryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PasswordRecover not implemented")
}
func (UnimplementedAuthServiceServer) Profile(context.Context, *ProfileRequest) (*ProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Profile not implemented")
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

func _AuthService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_RegisterConfirm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterConfirmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RegisterConfirm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_RegisterConfirm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RegisterConfirm(ctx, req.(*RegisterConfirmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Login(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_PasswordReset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).PasswordReset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_PasswordReset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).PasswordReset(ctx, req.(*PasswordResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_PasswordRecover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordRecoveryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).PasswordRecover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_PasswordRecover_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).PasswordRecover(ctx, req.(*PasswordRecoveryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Profile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Profile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Profile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Profile(ctx, req.(*ProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _AuthService_Register_Handler,
		},
		{
			MethodName: "RegisterConfirm",
			Handler:    _AuthService_RegisterConfirm_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AuthService_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _AuthService_Logout_Handler,
		},
		{
			MethodName: "PasswordReset",
			Handler:    _AuthService_PasswordReset_Handler,
		},
		{
			MethodName: "PasswordRecover",
			Handler:    _AuthService_PasswordRecover_Handler,
		},
		{
			MethodName: "Profile",
			Handler:    _AuthService_Profile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fightbettr.proto",
}

const (
	FightersService_SearchFightersCount_FullMethodName = "/FightersService/SearchFightersCount"
	FightersService_SearchFighters_FullMethodName      = "/FightersService/SearchFighters"
)

// FightersServiceClient is the client API for FightersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FightersServiceClient interface {
	SearchFightersCount(ctx context.Context, in *FightersRequest, opts ...grpc.CallOption) (*FightersCountResponse, error)
	SearchFighters(ctx context.Context, in *FightersRequest, opts ...grpc.CallOption) (*FightersResponse, error)
}

type fightersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFightersServiceClient(cc grpc.ClientConnInterface) FightersServiceClient {
	return &fightersServiceClient{cc}
}

func (c *fightersServiceClient) SearchFightersCount(ctx context.Context, in *FightersRequest, opts ...grpc.CallOption) (*FightersCountResponse, error) {
	out := new(FightersCountResponse)
	err := c.cc.Invoke(ctx, FightersService_SearchFightersCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fightersServiceClient) SearchFighters(ctx context.Context, in *FightersRequest, opts ...grpc.CallOption) (*FightersResponse, error) {
	out := new(FightersResponse)
	err := c.cc.Invoke(ctx, FightersService_SearchFighters_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FightersServiceServer is the server API for FightersService service.
// All implementations must embed UnimplementedFightersServiceServer
// for forward compatibility
type FightersServiceServer interface {
	SearchFightersCount(context.Context, *FightersRequest) (*FightersCountResponse, error)
	SearchFighters(context.Context, *FightersRequest) (*FightersResponse, error)
	mustEmbedUnimplementedFightersServiceServer()
}

// UnimplementedFightersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFightersServiceServer struct {
}

func (UnimplementedFightersServiceServer) SearchFightersCount(context.Context, *FightersRequest) (*FightersCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchFightersCount not implemented")
}
func (UnimplementedFightersServiceServer) SearchFighters(context.Context, *FightersRequest) (*FightersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchFighters not implemented")
}
func (UnimplementedFightersServiceServer) mustEmbedUnimplementedFightersServiceServer() {}

// UnsafeFightersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FightersServiceServer will
// result in compilation errors.
type UnsafeFightersServiceServer interface {
	mustEmbedUnimplementedFightersServiceServer()
}

func RegisterFightersServiceServer(s grpc.ServiceRegistrar, srv FightersServiceServer) {
	s.RegisterService(&FightersService_ServiceDesc, srv)
}

func _FightersService_SearchFightersCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FightersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FightersServiceServer).SearchFightersCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FightersService_SearchFightersCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FightersServiceServer).SearchFightersCount(ctx, req.(*FightersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FightersService_SearchFighters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FightersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FightersServiceServer).SearchFighters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FightersService_SearchFighters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FightersServiceServer).SearchFighters(ctx, req.(*FightersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FightersService_ServiceDesc is the grpc.ServiceDesc for FightersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FightersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FightersService",
	HandlerType: (*FightersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchFightersCount",
			Handler:    _FightersService_SearchFightersCount_Handler,
		},
		{
			MethodName: "SearchFighters",
			Handler:    _FightersService_SearchFighters_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fightbettr.proto",
}
