package services

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	mock_logger "projects/fb-server/pkg/logger/mocks"
	"projects/fb-server/pkg/model"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func getTestToken() (jwt.Token, error) {
	tokenId, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Unable to generate token id: %s", err)
	}

	token, err := jwt.NewBuilder().
		JwtID(tokenId.String()).
		Issuer("fb-fightbettr").
		Audience([]string{"localhost"}).
		IssuedAt(time.Now()).
		Subject("test").
		Expiration(time.Now().Add(5 * time.Second)).
		Build()

	if err != nil {
		return nil, err
	}

	return token, nil
}

func TestVerifyJWT(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockLogger := mock_logger.NewMockFbLogger(ctl)
	app := New(mockLogger, "TestApp")

	token, err := getTestToken()

	require.NoError(t, err)

	viper.Set("auth.jwt.cert", "../../hack/dev/certs/server-cert.pem")
	viper.Set("auth.jwt.key", "../../hack/dev/certs/server-key.pem")

	defer viper.Reset()

	mockLogger.EXPECT().Debugw(gomock.Any(), gomock.Any()).Times(1)

	app.loadJwtCerts()

	tests := []struct {
		name string
		test func()
	}{
		{
			name: "Correct payload",
			test: func() {
				alg := jwa.RS256
				payload, err := jwt.Sign(token, jwt.WithKey(alg, viper.Get("auth.jwt.signing_key")))
				require.NoError(t, err)
				result, err := app.verifyJWT(string(payload))

				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
		{
			name: "Bad payload",
			test: func() {
				mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).Times(1)

				_, err = app.verifyJWT("wrong data")

				assert.Error(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.test()
		})
	}

}

func TestIfLoggedIn(t *testing.T) {
	viper.Set("auth.jwt.cert", "../../hack/dev/certs/server-cert.pem")
	viper.Set("auth.jwt.key", "../../hack/dev/certs/server-key.pem")
	defer viper.Reset()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockLogger := mock_logger.NewMockFbLogger(ctl)
	app := New(mockLogger, "TestApp")
	alg := jwa.RS256
	cookieName := "session"

	mockLogger.EXPECT().Debugw(gomock.Any(), gomock.Any()).Times(1)
	app.loadJwtCerts()

	token, err := getTestToken()
	require.NoError(t, err)

	payload, err := jwt.Sign(token, jwt.WithKey(alg, viper.Get("auth.jwt.signing_key")))
	require.NoError(t, err)

	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		t.Log("Fake handler called")
	}

	tests := []struct {
		name string
		test func(req *http.Request)
	}{
		{
			name: "Cookie has bad access token",
			test: func(req *http.Request) {
				cookie := &http.Cookie{
					Name:  cookieName,
					Value: string(payload),
				}

				req.AddCookie(cookie)

				w := httptest.NewRecorder()

				app.IfLoggedIn(fakeHandler).ServeHTTP(w, req)

				assert.Equal(t, w.Code, http.StatusUnauthorized, "Expected status %d, but got %d", http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Bad cookie name",
			test: func(req *http.Request) {
				cookie := &http.Cookie{
					Name:  "bad name",
					Value: string(payload),
				}

				req.AddCookie(cookie)

				w := httptest.NewRecorder()

				mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).Times(1)
				app.IfLoggedIn(fakeHandler).ServeHTTP(w, req)

				assert.Equal(t, w.Code, http.StatusUnauthorized, "Expected status %d, but got %d", http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Bad cookie value",
			test: func(req *http.Request) {
				cookie := &http.Cookie{
					Name:  cookieName,
					Value: "bad value",
				}

				req.AddCookie(cookie)

				w := httptest.NewRecorder()

				mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
				app.IfLoggedIn(fakeHandler).ServeHTTP(w, req)

				assert.Equal(t, w.Code, http.StatusUnauthorized, "Expected status %d, but got %d", http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Good Token",
			test: func(req *http.Request) {
				curToken, err := getTestToken()
				require.NoError(t, err)

				if err := curToken.Set(string(model.ContextUserId), 1); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}

				curPayload, err := jwt.Sign(curToken, jwt.WithKey(alg, viper.Get("auth.jwt.signing_key")))
				require.NoError(t, err)

				cookie := &http.Cookie{
					Name:  cookieName,
					Value: string(curPayload),
				}

				req.AddCookie(cookie)

				w := httptest.NewRecorder()

				mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
				app.IfLoggedIn(fakeHandler).ServeHTTP(w, req)

				assert.Equal(t, w.Code, http.StatusOK, "Expected status %d, but got %d", http.StatusOK, w.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			tc.test(req)
		})
	}
}

func TestCheckIsAdmin(t *testing.T) {
	viper.Set("auth.jwt.cert", "../../hack/dev/certs/server-cert.pem")
	viper.Set("auth.jwt.key", "../../hack/dev/certs/server-key.pem")
	defer viper.Reset()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockLogger := mock_logger.NewMockFbLogger(ctl)
	app := New(mockLogger, "TestApp")
	alg := jwa.RS256
	cookieName := "session"

	mockLogger.EXPECT().Debugw(gomock.Any(), gomock.Any()).Times(1)
	app.loadJwtCerts()

	token, err := getTestToken()
	require.NoError(t, err)

	if err := token.Set(string(model.ContextUserId), 1); err != nil {
		log.Fatalf("Unable to set JWT token userRoles: %s", err)
	}

	tests := []struct {
		name string
		test func(req *http.Request)
	}{
		{
			name: "Action is not allowed",
			test: func(req *http.Request) {
				payload, err := jwt.Sign(token, jwt.WithKey(alg, viper.Get("auth.jwt.signing_key")))
				require.NoError(t, err)

				cookie := &http.Cookie{
					Name:  cookieName,
					Value: string(payload),
				}

				req.AddCookie(cookie)

				w := httptest.NewRecorder()

				fakeHandler := func(w http.ResponseWriter, r *http.Request) {
					t.Log("Fake handler called")
				}

				app.CheckIsAdmin(app.IfLoggedIn(fakeHandler)).ServeHTTP(w, req)

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, response["message"], "action is allowed only for admins", "Expected that method must be allowed only for admins")
			},
		},
		{
			name: "Action is allowed",
			test: func(req *http.Request) {
				if err := token.Set(string(model.ContextFlags), float64(1)); err != nil {
					log.Fatalf("Unable to set JWT token userRoles: %s", err)
				}

				payload, err := jwt.Sign(token, jwt.WithKey(alg, viper.Get("auth.jwt.signing_key")))
				require.NoError(t, err)
				cookie := &http.Cookie{
					Name:  cookieName,
					Value: string(payload),
				}

				req.AddCookie(cookie)

				w := httptest.NewRecorder()

				fakeHandler := func(w http.ResponseWriter, r *http.Request) {
					flagsValue := r.Context().Value(model.ContextJWTPointer)
					assert.NotNil(t, flagsValue, "Flag must be there")
				}

				app.CheckIsAdmin(fakeHandler).ServeHTTP(w, req)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			tc.test(req)
		})
	}
}
