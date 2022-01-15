package main

import (
	"context"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"option-dance/cmd/config"
	"option-dance/core"
	"option-dance/handler/api"
	"option-dance/version"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgFile string
	notify  bool
	port    string
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "od-engine",
		Short: "option dance match engine based on mixin trusted group",
		Run: func(cmd *cobra.Command, args []string) {
			rootEntry()
		},
	}
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate option dance database table struct",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := core.NewMigrator()
			err := migrator.Run()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "6028", "server bind port")
	rootCmd.PersistentFlags().BoolVar(&notify, "notify", false, "enable message notify")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug log console output")
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(SettlementCmd)
	rootCmd.AddCommand(MultisigCmd)
	rootCmd.AddCommand(GroupManageCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	config.InitConfig(cfgFile, debug)
	if err := config.ValidateEngineConfig(notify); err != nil {
		log.Panicln(err)
	}
}

func rootEntry() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := api.InitEngineServer()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return srv.Shutdown(ctx)
		}
	})
	zap.L().Info("engine version: " + version.Version.String() + " Source Distributions")
	g.Go(func() error {
		zap.L().Info("http server started at port " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		return nil
	})

	g.Go(func() error {
		exchange, err := InitExchange()
		if err != nil {
			log.Fatal(err)
		}
		return exchange.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		zap.L().Error("engine exited", zap.Error(err))
	}
}
