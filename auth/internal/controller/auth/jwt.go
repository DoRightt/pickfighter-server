package auth

import (
	"context"
	"time"

	authmodel "pickfighter.com/auth/pkg/model"
	logs "pickfighter.com/pkg/logger"
	"pickfighter.com/pkg/model"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/spf13/viper"
)

// createJWTToken generates a JWT token for the provided user credentials and authentication request.
// It includes user-specific claims such as user ID, user roles, and audience.
// The token is signed using RS256 algorithm with the configured signing key.
// Function returns resulting authentication token information (model.AuthenticateRequest),
// including token ID, access token, and expiration time.
func (c *Controller) createJWTToken(ctx context.Context, creds *authmodel.UserCredentials, req *authmodel.AuthenticateRequest) (*authmodel.AuthenticateResult, error) {
	u, err := c.repo.FindUser(ctx, &authmodel.UserRequest{
		UserId: creds.UserId,
	})
	if err != nil {
		logs.Errorf("Failed to get user: %s", err)
		return nil, err
	}

	logs.Debugf("Issuing JWT token for User [%d:%s:%s]", creds.UserId, creds.Email, req.Subject)

	tokenId, err := uuid.NewV4()
	if err != nil {
		logs.Errorf("Unable to generate token id: %s", err)
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
		logs.Errorf("Unable to build JWT token: %s", err)
		return nil, err
	}

	if err := t.Set(string(model.ContextUserId), u.UserId); err != nil {
		logs.Errorf("Unable to set JWT token userRoles: %s", err)
		return nil, err
	}

	if u.Flags > 0 {
		if err := t.Set(string(model.ContextFlags), u.Flags); err != nil {
			logs.Errorf("Unable to set JWT token private claim key: %s", err)
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
		logs.Errorf("failed to generate signed payload: %s\n", err)
		return nil, err
	}

	result := authmodel.AuthenticateResult{
		TokenId:        tokenId.String(),
		AccessToken:    string(payload),
		ExpirationTime: time.Now().Add(time.Duration(req.ExpiresIn) * time.Second),
	}

	return &result, nil
}
