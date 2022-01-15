package main

import (
	"context"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"option-dance/cmd/config"
	"option-dance/handler/api"
	"option-dance/service/message"
	"option-dance/version"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgFile string
	port    string
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "od-api",
		Short: "option dance api server",
		Run: func(cmd *cobra.Command, args []string) {
			rootEntry()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "6024", "api server bind port")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug log console output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	config.InitConfig(cfgFile, debug)
	err := config.ValidateApiConfig()
	if err != nil {
		log.Panicln(err)
	}
}

func rootEntry() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server, err := InitApiServer()
	if err != nil {
		zap.L().Fatal("init api server error", zap.Error(err))
	}
	handler := api.InitializeApiServer(server)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return srv.Shutdown(ctx)
		}
	})

	zap.L().Info("server version: " + version.Version.String() + " Source Distributions")
	g.Go(func() error {
		zap.L().Info("http server started at port " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		return nil
	})

	g.Go(func() error {
		zap.L().Info("blaze server started")
		go message.MixinMessageListener()
		return nil
	})

	g.Go(func() error {
		zap.L().Info("db market tracker service started")
		return server.DbMarketService.SyncMarket(ctx)
	})

	if err := g.Wait(); err != nil {
		zap.L().Error("api server exited", zap.Error(err))
	}
}
