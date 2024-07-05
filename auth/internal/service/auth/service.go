package service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"strings"

	grpchandler "fightbettr.com/auth/internal/handler/grpc"
	lg "fightbettr.com/auth/pkg/logger"
	"fightbettr.com/auth/pkg/version"
	"fightbettr.com/gen"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ErrAuthCertsPathRequired = fmt.Errorf("authentication certificates path is required")

type ApiService struct {
	ServiceName string
	Handler     *grpchandler.Handler
	Server      *grpc.Server
	Logger      lg.FbLogger
}

func New() ApiService {
	srv := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	logger := lg.GetSugared()

	return ApiService{
		ServiceName: version.Name,
		Server:      srv,
		Logger:      logger,
	}
}

func (s *ApiService) Init(h *grpchandler.Handler) error {
	s.Handler = h
	reflection.Register(s.Server)
	gen.RegisterAuthServiceServer(s.Server, s.Handler)

	if err := s.loadJwtCerts(); err != nil {
		s.Logger.Errorf("Unable to load JWT certificates: %s", err)
		return err
	}

	return nil
}

func (s *ApiService) Run() error {
	port := viper.GetString("http.port")
	srvAddr := viper.GetString("http.addr")
	if len(srvAddr) < 1 || !strings.Contains(srvAddr, ":") {
		return fmt.Errorf("'%s' service address not specified", s.ServiceName)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return err
	}

	s.Logger.Infof("Start listen '%s' http: %s", s.ServiceName, srvAddr)
	fmt.Printf("Server is listening at: %s\n", srvAddr)

	return s.Server.Serve(lis)
}

// loadJwtCerts loads the JWT certificates required for authentication from the specified paths.
// It expects paths to the X.509 certificate (certPath) and private key (keyPath) in the configuration.
// The loaded keypair is used for signing JWT tokens, and the public key is used for token verification.
// If the certificate or key cannot be loaded or parsed, an error is returned.
// The loaded keys are set in the configuration for later use in JWT signing and parsing.
func (h *ApiService) loadJwtCerts() error {
	certPath := viper.GetString("auth.jwt.cert")
	keyPath := viper.GetString("auth.jwt.key")

	hasRsaKeys := len(certPath) > 0 && len(keyPath) > 0

	if !hasRsaKeys {
		return ErrAuthCertsPathRequired
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		fmt.Println("ERR", err)
		h.Logger.Errorf("Unable to load client keypair: %s", err)
		return err
	}

	viper.Set("auth.jwt.signing_key", cert.PrivateKey)

	clientCert, err := os.ReadFile(certPath)
	if err != nil {
		h.Logger.Errorf("Unable to read key file bytes: %s", err)
		return err
	}

	block, _ := pem.Decode(clientCert)
	readCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		h.Logger.Errorf("Unable to parse certificate: %s", err)
		return err
	}

	viper.Set("auth.jwt.parse_key", readCert.PublicKey)

	h.Logger.Debugw("Loaded jwt certs",
		"cert_path", viper.GetString("auth.jwt.cert"),
		"key_path", viper.GetString("auth.jwt.key"),
	)

	return nil
}
