package auth

import (
	"context"
	"projects/fb-server/pkg/model"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

func (s *service) createJWTToken(ctx context.Context, creds *model.UserCredentials, req model.AuthenticateRequest) (*model.AuthenticateResult, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenId, err := uuid.NewV4()
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = creds.UserId
	claims["token_id"] = tokenId
	claims["code"] = "your_code"
	claims["access_token"] = "your_access_token"
	claims["expiration_time"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString([]byte("my_secret_key")) // TODO
	if err != nil {
		return nil, err
	}

	authResult := &model.AuthenticateResult{
		UserId:         creds.UserId,
		TokenId:        tokenId.String(),
		Code:           "your_code",
		AccessToken:    tokenString,
		ExpirationTime: time.Now().Add(time.Hour * 24 * 7),
	}

	return authResult, nil
}
