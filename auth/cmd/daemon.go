package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"pickfighter.com/auth/internal/controller/auth"
	grpchandler "pickfighter.com/auth/internal/handler/grpc"
	"pickfighter.com/auth/internal/repository/psql"
	service "pickfighter.com/auth/internal/service/auth"
	"pickfighter.com/pkg/discovery"
	"pickfighter.com/pkg/discovery/consul"
	logs "pickfighter.com/pkg/logger"
	"pickfighter.com/pkg/model"
	"pickfighter.com/pkg/sigx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var allowedApiRoutes = []string{
	model.AuthService,
}

var errEmptyApiRoute = fmt.Errorf("one of the api routes (%s) should be specified", strings.Join(allowedApiRoutes, ","))

func init() {
	rootCmd.AddCommand(serveCmd)
}

// serveCmd represents the serve command. It is used to run the gRPC server with specified API routes.
var serveCmd = &cobra.Command{
	Use:              "serve",
	Short:            "Run gRPC Server",
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	route := args[0]

	app := service.New()

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(app.ServiceName)
	if err := registry.Register(ctx, instanceID, app.ServiceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, app.ServiceName); err != nil {
				logs.Error("Failed to report healthy state", zap.Error(err))
			}

			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, app.ServiceName)

	repo, err := psql.New(ctx)
	if err != nil {
		logs.Errorf("Unable to start postgresql connection: %s", err)
		return
	}
	defer repo.GracefulShutdown()

	ctl := auth.New(repo)
	h := grpchandler.New(ctl)

	err = app.Init(h)
	if err != nil {
		log.Fatalf("error while app initialization: %s", err)
	}

	viper.Set("api.route", route)

	sigx.Listen(func(signal os.Signal) {
		time.AfterFunc(15*time.Second, func() {
			logs.Fatal("Failed to shutdown normally. Closed after 15 sec shutdown")
			cancel()

			os.Exit(1)
		})

		app.Server.GracefulStop()
	})

	if err := app.Run(); err != nil {
		logs.Fatal("app error: %s", err)
		app.Server.GracefulStop()
	}
}
