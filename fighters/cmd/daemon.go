package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"fightbettr.com/fb-server/pkg/sigx"
	"fightbettr.com/fighters/internal/controller/fighters"
	grpchandler "fightbettr.com/fighters/internal/handler/grpc"
	"fightbettr.com/fighters/internal/repository/psql"
	service "fightbettr.com/fighters/internal/service/fighters"
	"fightbettr.com/fighters/pkg/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allowedApiRoutes = []string{
	model.FightersService,
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	route := args[0]

	app := service.New()

	repo, err := psql.New(ctx, logger)
	if err != nil {
		logger.Errorf("Unable to start postgresql connection: %s", err)
		app.GracefulShutdown(ctx, err.Error())
	}
	ctl := fighters.New(repo)
	h := grpchandler.New(ctl)

	app.Init(h)

	viper.Set("api.route", route)

	sigx.Listen(func(signal os.Signal) {
		time.AfterFunc(15*time.Second, func() {
			logger.Fatal("Failed to shutdown normally. Closed after 15 sec shutdown")
		})
		cancel()

		app.GracefulShutdown(ctx, signal.String())
	})

	if err := app.Run(); err != nil {
		app.GracefulShutdown(ctx, err.Error())
	}
}
