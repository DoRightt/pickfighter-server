package services

import (
	"context"
	"fmt"
	"net/http"
	"projects/fb-server/pkg/httplib"
	"projects/fb-server/pkg/model"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"
)

// verifyJWT parses a raw JWT string and verifies its signature using the specified algorithm and public key.
// It returns a parsed JWT token on success, or an error if parsing or verification fails.
func (h *ApiHandler) verifyJWT(jwtRawValue string) (jwt.Token, error) {
	alg := jwa.RS256

	token, err := jwt.Parse([]byte(jwtRawValue), jwt.WithKey(alg, viper.Get("auth.jwt.parse_key")))
	if err != nil {
		h.Logger.Debugf("Failed to parse JWT token: %s", err)
		return nil, err
	}

	return token, nil
}

// IfLoggedIn is a middleware that checks if a user is logged in based on the provided JWT token.
// If the token is valid, it extracts user information such as user ID and claims, and adds them to the request context.
// If the token is invalid or missing, it returns an unauthorized response.
func (h *ApiHandler) IfLoggedIn(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie(httplib.CookieName)
		if err != nil {
			h.Logger.Debugf("access token not found: %s", err)
			httplib.ErrorResponseJSON(w, http.StatusUnauthorized, http.StatusUnauthorized,
				fmt.Errorf("unauthorized request: auth cookie or headers not found"))
			return
		}

		token, err := h.verifyJWT(cookie.Value)
		if err != nil {
			h.Logger.Debugf("Failed to parse JWT token: %s", err)
			httplib.ErrorResponseJSON(w, http.StatusUnauthorized, http.StatusUnauthorized,
				fmt.Errorf("unauthorized request: invalid token format"))
			return
		}

		userId, ok := token.Get(string(model.ContextUserId))
		if !ok {
			httplib.ErrorResponseJSON(w, http.StatusUnauthorized, http.StatusUnauthorized,
				fmt.Errorf("illegal token, user id must be specified"))
			return
		}

		if token.Expiration().Unix() < time.Now().Unix() {
			httplib.ErrorResponseJSON(w, http.StatusUnauthorized, http.StatusUnauthorized,
				fmt.Errorf("token expired"))
			return
		}

		uid, valid := userId.(float64)
		if valid {
			ctx = context.WithValue(ctx, model.ContextUserId, int32(uid))
		}
		
		// TODO admin claim?
		// rootClaim, onBoard := token.Get(string(model.ContextClaim))
		// if onBoard {
		// 	claim, fit := rootClaim.(string)
		// 	if fit {
		// 		ctx = context.WithValue(ctx, model.ContextNamespaceClaims, claim)
		// 	}
		// }

		ctx = context.WithValue(ctx, model.ContextJWTPointer, token)

		fn(w, r.WithContext(ctx))
	}
}

// CheckIsAdmin is a middleware that verifies if the user is logged in as an administrator.
// It checks the "admin" claim in the JWT (JSON Web Token) stored in the request cookie.
// If the user is an administrator, the request continues; otherwise, it responds with an error.
func (h *ApiHandler) CheckIsAdmin(next http.HandlerFunc) http.HandlerFunc {
	// TODO mb claim should set in createJWTToken method
	return h.IfLoggedIn(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie(httplib.CookieName)
		if err == nil {
			token, err := h.verifyJWT(cookie.Value)
			if err != nil {
				h.Logger.Debugf("Failed to parse JWT token: %s", err)
			} else {
				f, ok := token.Get(string(model.ContextFlags))
				flag, valid := f.(float64)
				if !ok || int(flag) != 1 {
					httplib.ErrorResponseJSON(w, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed,
						fmt.Errorf("action is allowed only for admins"))
					return
				}

				if valid {
					ctx = context.WithValue(ctx, model.ContextFlags, int(flag))
				}

				ctx = context.WithValue(ctx, model.ContextJWTPointer, token)
			}

			r = r.WithContext(ctx)
		}

		next(w, r.WithContext(ctx))
	})
}
