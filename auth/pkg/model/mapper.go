package model

import (
	"pickfighter.com/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RegisterRequestFromProto(p *gen.RegisterRequest) *RegisterRequest {
	return &RegisterRequest{
		Name:     p.Name,
		Email:    p.Email,
		Password: p.Password,
		Token:    p.Token,
		TermsOk:  p.TermsOk,
	}
}

func RegisterRequestToProto(req *RegisterRequest) *gen.RegisterRequest {
	return &gen.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Token:    req.Token,
		TermsOk:  req.TermsOk,
	}
}

func AuthenticateRequestFromProto(p *gen.AuthenticateRequest) *AuthenticateRequest {
	return &AuthenticateRequest{
		Email:    p.Email,
		Password: p.Password,

		RememberMe: p.RememberMe,
		UserAgent:  p.UserAgent,
		IpAddress:  p.IpAddress,

		Subject:   p.Subject,
		ExpiresIn: p.ExpiresIn,
		Audience:  p.Audience,

		Method: int(p.Method),
	}
}

func AuthenticateRequestToProto(req *AuthenticateRequest) *gen.AuthenticateRequest {
	return &gen.AuthenticateRequest{
		Email:    req.Email,
		Password: req.Password,

		RememberMe: req.RememberMe,
		UserAgent:  req.UserAgent,
		IpAddress:  req.IpAddress,

		Subject:   req.Subject,
		ExpiresIn: req.ExpiresIn,
		Audience:  req.Audience,

		Method: int32(req.Method),
	}
}

func AuthenticateResultFromProto(p *gen.AuthenticateResponse) *AuthenticateResult {
	return &AuthenticateResult{
		// UserId:         p.UserId,
		// Code:           p.Code,
		TokenId:        string(p.TokenId),
		AccessToken:    p.AccessToken,
		ExpirationTime: p.ExpirationTime.AsTime(),
	}
}

func AuthenticateResultToProto(req *AuthenticateResult) *gen.AuthenticateResponse {
	return &gen.AuthenticateResponse{
		// UserId:         p.UserId,
		// Code:           p.Code,
		TokenId:        req.TokenId,
		AccessToken:    req.AccessToken,
		ExpirationTime: timestamppb.New(req.ExpirationTime),
	}
}

func PasswordRecoveryRequestFromProto(p *gen.PasswordRecoveryRequest) *RecoverPasswordRequest {
	return &RecoverPasswordRequest{
		Token:           p.Token,
		Password:        p.Password,
		ConfirmPassword: p.ConfirmPassword,
	}
}

func PasswordRecoveryRequestToProto(req *RecoverPasswordRequest) *gen.PasswordRecoveryRequest {
	return &gen.PasswordRecoveryRequest{
		Token:           req.Token,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
}

func UserFromProto(p *gen.ProfileResponse) *User {
	return &User{
		UserId:    p.User.UserId,
		Name:      p.User.Name,
		Email:     p.User.Email,
		Rank:      p.User.Rank,
		Claim:     p.User.Claim,
		Roles:     p.User.Roles,
		Flags:     p.User.Flags,
		CreatedAt: p.User.CreatedAt,
		UpdatedAt: p.User.UpdatedAt,
	}
}

func UserToProto(u *User) *gen.ProfileResponse {
	return &gen.ProfileResponse{
		User: &gen.User{
			UserId:    u.UserId,
			Name:      u.Name,
			Email:     u.Email,
			Rank:      u.Rank,
			Claim:     u.Claim,
			Roles:     u.Roles,
			Flags:     u.Flags,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}
}
