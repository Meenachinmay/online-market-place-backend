package bootstrap

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"soda-interview/foundation/config"
	"soda-interview/foundation/database/postgres"
	"soda-interview/foundation/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RegisterFn func(log *logger.Logger, db *pgxpool.Pool, grpcServer *grpc.Server)

// Run initializes the system infrastructure and starts the gRPC server.
// It delegates the specific service registration to the register callback.
func Run(register RegisterFn) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(os.Stdout, cfg.Logging.Level)
	log.Info("Starting service", "name", cfg.App.Name, "env", cfg.App.Environment)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := postgres.New(ctx, cfg.GetDatabaseDSN())
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := postgres.Migrate(context.Background(), dbPool, "business/data/schema/migrations", log.NewStdLogger()); err != nil {
		log.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", cfg.GetGRPCAddress())
	if err != nil {
		log.Error("Failed to listen", "address", cfg.GetGRPCAddress(), "error", err)
		os.Exit(1)
	}

	var serverOpts []grpc.ServerOption
	gRPCServer := grpc.NewServer(serverOpts...)

	// Register Health Service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	log.Info("gRPC health service registered")

	register(log, dbPool, gRPCServer)

	if !cfg.IsProduction() {
		reflection.Register(gRPCServer)
		log.Info("gRPC reflection enabled")
	}

	go func() {
		log.Info("gRPC server starting", "address", cfg.GetGRPCAddress())
		if err := gRPCServer.Serve(lis); err != nil {
			log.Error("Failed to serve gRPC", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	_, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	gRPCServer.GracefulStop()
	log.Info("Server stopped")
}