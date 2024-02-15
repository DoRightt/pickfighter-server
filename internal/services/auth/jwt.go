package auth

import (
	"context"
	"projects/fb-server/pkg/model"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/spf13/viper"
)

func (s *service) createJWTToken(ctx context.Context, creds *model.UserCredentials, req model.AuthenticateRequest) (*model.AuthenticateResult, error) {
	u, err := s.Repo.FindUser(ctx, &model.UserRequest{
		UserId: creds.UserId,
	})
	if err != nil {
		s.Logger.Errorf("Failed to get user: %s", err)
		return nil, err
	}

	s.Logger.Debugf("Issuing JWT token for User [%d:%s:%s]", creds.UserId, creds.Email, req.Subject)

	tokenId, err := uuid.NewV4()
	if err != nil {
		s.Logger.Errorf("Unable to generate token id: %s", err)
		return nil, err
	}

	req.Audience = append(req.Audience, viper.GetString("domain"))

	now := time.Now()
	subject := crypto.Keccak256([]byte(req.IpAddress + ":" + req.UserAgent))

	t, err := jwt.NewBuilder().
		JwtID(tokenId.String()).
		Issuer(viper.GetString("auth.issuer")).
		Audience(req.Audience).
		IssuedAt(now).
		Subject(hexutil.Encode(subject)).
		Expiration(now.Add(time.Duration(req.ExpiresIn) * time.Second)).
		Build()
	if err != nil {
		s.Logger.Errorf("Unable to build JWT token: %s", err)
		return nil, err
	}
	
	if err := t.Set(model.ContextUserId, u.UserId); err != nil {
		s.Logger.Errorf("Unable to set JWT token userRoles: %s", err)
		return nil, err
	}

	if u.Flags > 0 {
		if err := t.Set(model.ContextFlags, u.Flags); err != nil {
			s.Logger.Errorf("Unable to set JWT token private claim key: %s", err)
			return nil, err
		}
	}

	// buf, err := json.MarshalIndent(t, "", "  ")
	// if err != nil {
	// 	s.Logger.Errorf("Failed to generate JWT Token JSON: %s", err)
	// 	return nil, err
	// }

	alg := jwa.RS256
	payload, err := jwt.Sign(t, jwt.WithKey(alg, viper.Get("auth.jwt.signing_key"))) // TODO
	if err != nil {
		s.Logger.Errorf("failed to generate signed payload: %s\n", err)
		return nil, err
	}

	result := model.AuthenticateResult{
		TokenId:        tokenId.String(),
		AccessToken:    string(payload),
		ExpirationTime: time.Now().Add(time.Duration(req.ExpiresIn) * time.Second),
	}

	return &result, nil
}
