package grpc

import (
	"context"

	authmodel "fightbettr.com/auth/pkg/model"
	"fightbettr.com/gen"
	"fightbettr.com/internal/grpcutil"
	"fightbettr.com/pkg/discovery"
	"fightbettr.com/pkg/model"
)

// Gateway defines an gRPC gateway for a auth service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a auth service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Register registers a new user via the auth-service.
// It establishes a gRPC connection, sends the registration request,
// and returns the user's credentials if successful.
func (g *Gateway) Register(ctx context.Context, req *authmodel.RegisterRequest) (*authmodel.UserCredentials, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "auth-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	regReq := authmodel.RegisterRequestToProto(req)
	resp, err := client.Register(ctx, regReq)
	if err != nil {
		return nil, err
	}

	credentials := &authmodel.UserCredentials{UserId: resp.Id}

	return credentials, nil
}

func (g *Gateway) ConfirmRegistration(ctx context.Context, token string) (bool, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "auth-service", g.registry)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	confReq := &gen.RegisterConfirmRequest{Token: token}
	_, err = client.RegisterConfirm(ctx, confReq)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Login authenticates a user via the auth-service.
// It establishes a gRPC connection, sends the authentication request,
// and returns the authentication result if successful.
func (g *Gateway) Login(ctx context.Context, req *authmodel.AuthenticateRequest) (*authmodel.AuthenticateResult, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "auth-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	loginReq := authmodel.AuthenticateRequestToProto(req)
	resp, err := client.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	token := authmodel.AuthenticateResultFromProto(resp)

	return token, nil
}

// ResetPassword initiates a password reset process via the auth-service.
// It establishes a gRPC connection, sends the password reset request,
// and returns true if the request was successfully processed.
func (g *Gateway) ResetPassword(ctx context.Context, req *authmodel.ResetPasswordRequest) (bool, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "auth-service", g.registry)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	passResetReq := &gen.PasswordResetRequest{Email: req.Email}
	_, err = client.PasswordReset(ctx, passResetReq)
	if err != nil {
		return false, err
	}

	return true, nil
}

// PasswordRecover initiates the password recovery process via the auth-service.
// It establishes a gRPC connection, sends the password recovery request,
// and returns true if the request was successfully processed.
func (g *Gateway) PasswordRecover(ctx context.Context, req *authmodel.RecoverPasswordRequest) (bool, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "auth-service", g.registry)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	passRecoverReq := authmodel.PasswordRecoveryRequestToProto(req)
	_, err = client.PasswordRecover(ctx, passRecoverReq)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetCurrentUser retrieves the current authenticated user's profile via the auth-service.
// It establishes a gRPC connection, sends a profile request,
// and returns the user's profile if successfully retrieved.
func (g *Gateway) GetCurrentUser(ctx context.Context) (*authmodel.User, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "auth-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	userId := ctx.Value(model.ContextUserId).(int32)
	profileReq := &gen.ProfileRequest{UserId: userId}

	resp, err := client.Profile(ctx, profileReq)
	if err != nil {
		return nil, err
	}

	user := authmodel.UserFromProto(resp)

	return user, nil
}
