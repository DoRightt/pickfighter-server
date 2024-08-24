package service

import (
	"fmt"
	"net"
	"strings"

	grpchandler "pickfighter.com/fighters/internal/handler/grpc"
	logs "pickfighter.com/pkg/logger"
	"pickfighter.com/fighters/pkg/version"
	"pickfighter.com/gen"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiService struct {
	ServiceName string
	Handler     *grpchandler.Handler
	Server      *grpc.Server
}

func New() ApiService {
	srv := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))

	return ApiService{
		ServiceName: version.Name,
		Server:      srv,
	}
}

func (s *ApiService) Init(h *grpchandler.Handler) {
	s.Handler = h
	reflection.Register(s.Server)
	gen.RegisterFightersServiceServer(s.Server, s.Handler)
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

	logs.Infof("Start listen '%s' http: %s", s.ServiceName, srvAddr)
	fmt.Printf("Server is listening at: %s\n", srvAddr)

	return s.Server.Serve(lis)
}
