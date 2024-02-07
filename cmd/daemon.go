package cmd

import (
	"context"
	"fmt"
	"os"
	"projects/fb-server/internal/services"
	"projects/fb-server/internal/services/auth"
	"projects/fb-server/internal/services/common"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/sigx"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allowedApiRoutes = []string{
	model.AuthService,
	model.CommonService,
}

var errEmptyApiRoute = fmt.Errorf("one of the api routes (%s) should be specified", strings.Join(allowedApiRoutes, ","))

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("addr", ":9090", "Specify service routes to serve")
	bindViperFlag(serveCmd, "http.addr", "addr")

	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
}

var serveCmd = &cobra.Command{
	Use:              "serve",
	Short:            "Run HTTP Server",
	Long:             ``,
	TraverseChildren: true,
	Args: func(cmd *cobra.Command, args []string) error {
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		route := args[0]

		app := services.New(logger, route)

		sigx.Listen(func(signal os.Signal) {
			time.AfterFunc(15*time.Second, func() {
				app.Logger.Fatal("Failed to shutdown normally. Closed after 15 sec shutdown")
			})

			app.GracefulShutdown(ctx, signal.String())
		})

		if err := app.Init(ctx); err != nil {
			app.GracefulShutdown(ctx, err.Error())
		}

		viper.Set("api.route", route)
		switch route {
		case model.AuthService:
			app.AddService(model.AuthService, auth.New(app))
			break
		case model.CommonService:
			app.AddService(model.CommonService, common.New(app))
		default:
			app.GracefulShutdown(ctx, "invalid service route")
		}

		if err := app.Run(ctx); err != nil {
			app.GracefulShutdown(ctx, err.Error())
		}
	},
}

var stopCmd = &cobra.Command{}

var statusCmd = &cobra.Command{}
