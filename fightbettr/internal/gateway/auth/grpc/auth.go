package grpc

import (
	"context"

	authmodel "fightbettr.com/auth/pkg/model"
	"fightbettr.com/gen"
	"fightbettr.com/internal/grpcutil"
	"fightbettr.com/pkg/discovery"
)

// Gateway defines an gRPC gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

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

func (g *Gateway) GetCurrentUser(ctx context.Context) (*authmodel.User, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "auth-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	profileReq := &gen.ProfileRequest{} // TBD
	resp, err := client.Profile(ctx, profileReq)
	if err != nil {
		return nil, err
	}

	user := authmodel.UserFromProto(resp)

	return user, nil
}
