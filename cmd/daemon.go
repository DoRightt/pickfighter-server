package cmd

import (
	"context"
	"fmt"
	"os"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/sigx"
	"projects/fb-server/services"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var allowedApiRoutes = []string{
	model.AuthService,
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
		fmt.Println(len(args))
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

			fmt.Println("TEST sig", signal)

			app.GracefulShutdown(ctx, signal.String())
		})

		if err := app.Init(ctx); err != nil {
			app.GracefulShutdown(ctx, err.Error())
		}

		if err := app.Run(ctx); err != nil {
			app.GracefulShutdown(ctx, err.Error())
		}
	},
}

var stopCmd = &cobra.Command{}

var statusCmd = &cobra.Command{}
