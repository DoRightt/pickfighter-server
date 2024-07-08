package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	lg "fightbettr.com/fightbettr/pkg/logger"
	"fightbettr.com/fightbettr/pkg/version"
	"github.com/spf13/viper"
)

var ErrAuthCertsPathRequired = fmt.Errorf("authentication certificates path is required")

type HttpHandler interface {
	RunHTTPServer(ctx context.Context) error
}

type ApiService struct {
	ServiceName string
	Handler     HttpHandler
	Logger      lg.FbLogger
}

// New gets logger and returns new instance of ApiService
func New(h HttpHandler) ApiService {
	logger := lg.GetSugared()

	return ApiService{
		ServiceName: version.Name,
		Handler:     h,
		Logger:      logger,
	}
}

// Run starts the API service's HTTP server.
func (s *ApiService) Run(ctx context.Context) error {
	if err := s.loadJwtCerts(); err != nil {
		s.Logger.Errorf("Unable to load JWT certificates: %s", err)
		return err
	}

	return s.Handler.RunHTTPServer(ctx)
}

// loadJwtCerts loads the JWT certificates required for authentication from the specified paths.
// It expects paths to the X.509 certificate (certPath) and private key (keyPath) in the configuration.
// The loaded keypair is used for signing JWT tokens, and the public key is used for token verification.
// If the certificate or key cannot be loaded or parsed, an error is returned.
// The loaded keys are set in the configuration for later use in JWT signing and parsing.
func (s *ApiService) loadJwtCerts() error {
	certPath := viper.GetString("auth.jwt.cert")
	keyPath := viper.GetString("auth.jwt.key")

	hasRsaKeys := len(certPath) > 0 && len(keyPath) > 0

	if !hasRsaKeys {
		return ErrAuthCertsPathRequired
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		fmt.Println("ERR", err)
		s.Logger.Errorf("Unable to load client keypair: %s", err)
		return err
	}

	viper.Set("auth.jwt.signing_key", cert.PrivateKey)

	clientCert, err := os.ReadFile(certPath)
	if err != nil {
		s.Logger.Errorf("Unable to read key file bytes: %s", err)
		return err
	}

	block, _ := pem.Decode(clientCert)
	readCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		s.Logger.Errorf("Unable to parse certificate: %s", err)
		return err
	}

	viper.Set("auth.jwt.parse_key", readCert.PublicKey)

	s.Logger.Debugw("Loaded jwt certs",
		"cert_path", viper.GetString("auth.jwt.cert"),
		"key_path", viper.GetString("auth.jwt.key"),
	)

	return nil
}

// GracefulShutdown logs the received signal and exits the service.
func (s *ApiService) GracefulShutdown(ctx context.Context, sig string) {
	s.Logger.Warnf("Graceful shutdown. Signal received: %s", sig)

	os.Exit(0)
}
