package service

import (
	"fmt"
	"net"
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

func (s *ApiService) Init(h *grpchandler.Handler) {
	s.Handler = h
	reflection.Register(s.Server)
	gen.RegisterAuthServiceServer(s.Server, s.Handler)
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