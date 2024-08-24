package grpc

import (
	"context"

	"pickfighter.com/auth/internal/controller/auth"
	"pickfighter.com/auth/pkg/model"
	"pickfighter.com/gen"
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

// Register handles user registration by converting the gRPC request to internal format,
// delegating to the controller, and returning the registered user's ID or an error.
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

// RegisterConfirm handles the gRPC request to confirm user registration.
// It verifies the token from the request, delegates to the controller for confirmation,
// and returns success or an error if confirmation fails.
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

// Login handles the gRPC request to authenticate a user.
// It converts the protobuf request to internal format, delegates to the controller,
// and returns the authentication response or an error if authentication fails.
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

// PasswordReset handles the gRPC request to reset a user's password.
// It verifies the request, converts it to internal format, delegates to the controller,
// and returns success or an error if the password reset fails.
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

// PasswordRecover handles the gRPC request to initiate the password recovery process.
// It validates the request, converts it to internal format, delegates to the controller
// for further processing, and returns success or an error if password recovery fails.
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

// Profile handles the gRPC request to fetch the user profile based on the current user ID in the context.
// It verifies the presence of the user ID in the context, retrieves the user profile through the controller,
// and returns the profile information or an error if fetching the profile fails.
func (h *Handler) Profile(ctx context.Context, req *gen.ProfileRequest) (*gen.ProfileResponse, error) {
	currentUserId := req.UserId
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
