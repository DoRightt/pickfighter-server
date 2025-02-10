package service

import (
	"fmt"
	"net"

	grpchandler "github.com/DoRightt/pickfighter-server/auth/internal/handler/grpc"
	"github.com/DoRightt/pickfighter-server/auth/pkg/version"
	"github.com/DoRightt/pickfighter-server/gen"
	logs "github.com/DoRightt/pickfighter-server/pkg/logger"
	"github.com/DoRightt/pickfighter-server/pkg/utils"
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
}

func New() ApiService {
	srv := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))

	return ApiService{
		ServiceName: version.Name,
		Server:      srv,
	}
}

func (s *ApiService) Init(h *grpchandler.Handler) error {
	s.Handler = h
	reflection.Register(s.Server)
	gen.RegisterAuthServiceServer(s.Server, s.Handler)

	if err := utils.LoadJwtCerts(); err != nil {
		logs.Errorf("Unable to load JWT certificates: %s", err)
		return err
	}

	return nil
}

func (s *ApiService) Run() error {
	port := viper.GetString("http.port")
	srvAddr := viper.GetString("http.addr")

	if len(srvAddr) < 1 || len(port) < 1 {
		return fmt.Errorf("'%s' service address not specified", s.ServiceName)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%v", srvAddr, port))
	if err != nil {
		return err
	}

	logs.Infof("Start listen '%s' http: %s", s.ServiceName, srvAddr)
	fmt.Printf("Server is listening at: %s\n", fmt.Sprintf("%s:%v", srvAddr, port))

	return s.Server.Serve(lis)
}
