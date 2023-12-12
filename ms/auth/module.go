package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pb "github.com/jessicatarra/greenlight/api/proto"
	"github.com/jessicatarra/greenlight/internal/config"
	appl "github.com/jessicatarra/greenlight/ms/auth/internal/application"
	_grpc "github.com/jessicatarra/greenlight/ms/auth/internal/infrastructure/grpc"
	_http "github.com/jessicatarra/greenlight/ms/auth/internal/infrastructure/http"
	repo "github.com/jessicatarra/greenlight/ms/auth/internal/infrastructure/repositories"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	defaultIdleTimeout  = time.Minute
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

type module struct {
	grpc   *grpc.Server
	server *http.Server
	logger *slog.Logger
}

func (m module) Start(wg *sync.WaitGroup) {
	wg.Add(2)
	go func() {
		defer wg.Done()
		m.logger.Info("Starting Auth Module server", slog.Group("server", "addr", m.server.Addr))

		err := m.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			m.logger.Info("Auth module encountered an error")

			os.Exit(1)
		}
		m.logger.Info("Stopped auth Module server", slog.Group("server", "addr", m.server.Addr))

	}()

	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		m.logger.Info("Starting Auth Module GRPC server", slog.Group("server", "addr", lis.Addr()))
		if err := m.grpc.Serve(lis); err != nil {
			m.logger.Info("Auth module encountered an error")
			os.Exit(1)
		}
		m.logger.Info("Stopped auth Module GRPC server", slog.Group("server", "addr", lis.Addr()))
	}()
}

func (m module) Shutdown(ctx context.Context, cancel func()) {
	defer cancel()

	m.grpc.GracefulStop()
	err := m.server.Shutdown(ctx)
	if err != nil {
		return
	}

}

func NewModule(db *sql.DB, cfg config.Config, wg *sync.WaitGroup, logger *slog.Logger) *module {
	userRepo := repo.NewUserRepo(db)
	tokenRepo := repo.NewTokenRepo(db)
	permissionRepo := repo.NewPermissionRepo(db)
	appl := appl.NewAppl(userRepo, tokenRepo, permissionRepo, wg, cfg)
	api := _http.NewService(appl, cfg, logger)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthGRPCServiceServer(grpcServer, _grpc.NewGRPCServer(appl))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8082),
		Handler:      api.Routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &module{grpc: grpcServer, server: srv, logger: logger}
}
