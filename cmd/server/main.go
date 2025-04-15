package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arefev/mtrcstore/internal/proto"
	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/service"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	config, err := NewConfig(args)
	if err != nil {
		return fmt.Errorf("main config init failed: %w", err)
	}

	cLog, err := logger.Build(config.LogLevel)
	if err != nil {
		return fmt.Errorf("logger init failed: %w", err)
	}

	storage, err := initStorage(&config, cLog)
	if err != nil {
		return fmt.Errorf("main run failed: %w", err)
	}

	defer func() {
		if err := storage.Close(); err != nil {
			cLog.Error("storage close failed: %w", zap.Error(err))
		}
	}()

	switch {
	case config.GRPCAddress != "":
		return runGRPC(ctx, storage, &config, cLog)
	default:
		return runServer(ctx, storage, &config, cLog)
	}
}

func runGRPC(ctx context.Context, storage repository.Storage, c *Config, l *zap.Logger) error {
	listen, err := net.Listen("tcp", c.GRPCAddress)
	if err != nil {
		return fmt.Errorf("runGRPC Listen failed: %w", err)
	}

	s := grpc.NewServer()
	proto.RegisterMetricsServer(s, &service.GRPCServer{
		Storage: storage,
	})

	go func() {
		<-ctx.Done()
		l.Info("GRPC stopped")
		s.Stop()
	}()

	l.Info(
		"GRPC running",
		zap.String("address", c.GRPCAddress),
		zap.String("log level", c.LogLevel),
	)

	if err := s.Serve(listen); err != nil {
		return fmt.Errorf("runGRPC Serve failed: %w", err)
	}

	return nil
}

func runServer(ctx context.Context, storage repository.Storage, c *Config, l *zap.Logger) error {
	metricHandlers := handler.NewMetricHandlers(storage, l)
	r := server.InitRouter(metricHandlers, l, c.TrustedSubnet, c.SecretKey, c.CryptoKey)

	g, gCtx := errgroup.WithContext(ctx)
	serv := http.Server{
		Addr:    c.Address,
		Handler: r,
		BaseContext: func(_ net.Listener) context.Context {
			return gCtx
		},
	}

	g.Go(func() error {
		<-ctx.Done()
		l.Info("Server stopped")
		return serv.Shutdown(ctx)
	})

	l.Info(
		"Server running",
		zap.String("address", c.Address),
		zap.String("log level", c.LogLevel),
	)

	g.Go(serv.ListenAndServe)

	if err := g.Wait(); err != nil {
		return fmt.Errorf("exit reason: %w", err)
	}

	return nil
}

func initStorage(config *Config, cLog *zap.Logger) (repository.Storage, error) {
	var storage repository.Storage
	var err error

	switch {
	case len(config.DatabaseDSN) > 0:
		storage, err = repository.NewDatabaseRep(config.DatabaseDSN, cLog)
		if err != nil {
			err = fmt.Errorf("repository init failed: %w", err)
		}
	case len(config.FileStoragePath) > 0:
		storage = repository.
			NewFile(config.StoreInterval, config.FileStoragePath, config.Restore, cLog).
			WorkerRun()
	default:
		storage = repository.NewMemory()
	}

	return storage, err
}
