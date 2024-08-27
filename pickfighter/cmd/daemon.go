package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"pickfighter.com/pickfighter/internal/controller/pickfighter"
	authgateway "pickfighter.com/pickfighter/internal/gateway/auth/grpc"
	eventgateway "pickfighter.com/pickfighter/internal/gateway/events/grpc"
	fightersgateway "pickfighter.com/pickfighter/internal/gateway/fighters/grpc"
	httphandler "pickfighter.com/pickfighter/internal/handler/http"
	service "pickfighter.com/pickfighter/internal/service/pickfighter"
	logs "pickfighter.com/pkg/logger"
	"pickfighter.com/pickfighter/pkg/version"
	"pickfighter.com/pkg/discovery"
	"pickfighter.com/pkg/discovery/consul"
	"pickfighter.com/pkg/model"
	"pickfighter.com/pkg/sigx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var allowedApiRoutes = []string{
	model.GatewayService,
}

var errEmptyApiRoute = fmt.Errorf("one of the api routes (%s) should be specified", strings.Join(allowedApiRoutes, ","))

func init() {
	rootCmd.AddCommand(serveCmd)
}

// serveCmd represents the serve command. It is used to run the HTTP server with specified API routes.
var serveCmd = &cobra.Command{
	Use:              "serve",
	Short:            "Run HTTP Server",
	Long:             ``,
	TraverseChildren: true,
	Args:             validateServerArgs,
	Run:              runServe,
}

// validateServerArgs is a function used to validate the arguments passed to the serve command.
// It checks if a single API route is provided and if it is valid.
func validateServerArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errEmptyApiRoute
	}

	var ok bool
	for i := range allowedApiRoutes {
		if allowedApiRoutes[i] == args[0] {
			ok = true
			break
		}
	}

	if !ok {
		return fmt.Errorf("allowed routes are: %s", strings.Join(allowedApiRoutes, ", "))
	}

	return nil
}

// runServe is the main function executed when the serve command is run.
// It initializes the application, sets up service and runs the HTTP server.
func runServe(cmd *cobra.Command, args []string) {
	port := viper.GetInt("http.port")
	serviceName := version.Name
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	route := args[0]

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				logs.Error("Failed to report healthy state", zap.Error(err))
			}

			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	authGateway := authgateway.New(registry)
	eventGateway := eventgateway.New(registry)
	fightersGateway := fightersgateway.New(registry)
	ctl := pickfighter.New(authGateway, eventGateway, fightersGateway)
	h := httphandler.New(ctl)
	app := service.New(h)

	viper.Set("api.route", route)

	sigx.Listen(func(signal os.Signal) {
		time.AfterFunc(15*time.Second, func() {
			logs.Fatal("Failed to shutdown normally. Closed after 15 sec shutdown")
			cancel()

			os.Exit(1)
		})

		app.GracefulShutdown(ctx, signal.String())
	})

	if err := app.Run(ctx); err != nil {
		app.GracefulShutdown(ctx, err.Error())
	}
}
