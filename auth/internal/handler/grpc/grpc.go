package grpc

import (
	"context"

	"fightbettr.com/auth/internal/controller/auth"
	"fightbettr.com/auth/pkg/model"
	"fightbettr.com/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a Auth gRPC handler.
type Handler struct {
	gen.UnimplementedAuthServiceServer
	ctrl *auth.Controller
}

// New creates a new Auth gRPC handler.
func New(ctrl *auth.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	regReq := model.RegisterRequestFromProto(req)
	v, err := h.ctrl.Register(ctx, regReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.RegisterResponse{Id: v}, nil
}

func (h *Handler) RegisterConfirm(ctx context.Context, req *gen.RegisterConfirmRequest) (*gen.RegisterConfirmResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	confReq := &model.UserCredentialsRequest{Token: req.Token}
	_, err := h.ctrl.RegisterConfirm(ctx, confReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.RegisterConfirmResponse{}, nil
}

func (h *Handler) Login(ctx context.Context, req *gen.AuthenticateRequest) (*gen.AuthenticateResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	loginReq := model.AuthenticateRequestFromProto(req)
	resp, err := h.ctrl.Login(ctx, loginReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.AuthenticateResultToProto(resp), nil
}

func (h *Handler) PasswordReset(ctx context.Context, req *gen.PasswordResetRequest) (*gen.PasswordResetResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	resetPasswordReq := &model.ResetPasswordRequest{Email: req.Email}
	_, err := h.ctrl.PasswordReset(ctx, resetPasswordReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.PasswordResetResponse{}, nil
}

func (h *Handler) PasswordRecover(ctx context.Context, req *gen.PasswordRecoveryRequest) (*gen.PasswordRecoveryResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	passRecoverReq := model.PasswordRecoveryRequestFromProto(req)
	_, err := h.ctrl.PasswordRecover(ctx, passRecoverReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.PasswordRecoveryResponse{}, nil
}

// TODO req TBD
func (h *Handler) Profile(ctx context.Context, req *gen.ProfileRequest) (*gen.ProfileResponse, error) {
	currentUserId := ctx.Value(model.ContextUserId).(int32)
	if currentUserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "request has no id")
	}

	profileReq := &model.UserRequest{UserId: currentUserId}
	user, err := h.ctrl.Profile(ctx, profileReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.UserToProto(user), nil
}

func (h *Handler) GracefulShutdown(ctx context.Context, sig string) {
	h.ctrl.GracefulShutdown(ctx, sig)
}
