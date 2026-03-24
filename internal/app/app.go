package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/slavatrudu/auth/internal/config"
	_ "github.com/slavatrudu/auth/internal/migrations"
	"github.com/slavatrudu/auth/internal/repository"
	"github.com/slavatrudu/auth/internal/server"
	"github.com/slavatrudu/auth/internal/service"
	authpb "github.com/slavatrudu/contracts/auth/go"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger

	repo       *repository.Repository
	auth       *service.AuthService
	authServer *server.Server
	grpcServer *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	// Инициализируем gRPC сервер
	authServer, err := a.getAuthServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth server: %w", err)
	}

	// Создаем gRPC сервер
	a.grpcServer = getGRPCServer(authServer)

	listenAddr := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		a.logger.Fatal().Err(err).Msg("failed to listen")
		return err
	}

	a.logger.Info().Msg("gRPC server listening")

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- a.grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		a.grpcServer.GracefulStop()
		return ctx.Err()
	case err := <-serveErrCh:
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to serve")
		}
		return err
	}
}

func (a *App) getRepository(ctx context.Context) (*repository.Repository, error) {
	if a.repo == nil {
		if err := a.runMigrations(ctx); err != nil {
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}
		db, err := gorm.Open(postgres.Open(a.cfg.DbDsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			return nil, fmt.Errorf("gorm init failed: %w", err)
		}
		a.repo = repository.NewRepository(db, a.logger)
	}
	return a.repo, nil
}

func (a *App) runMigrations(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to select migrations dialect: %w", err)
	}

	dbGoose, err := sql.Open("postgres", a.cfg.DbDsn)
	if err != nil {
		return fmt.Errorf("failed to create sql connection: %w", err)
	}

	if err := goose.UpContext(ctx, dbGoose, "internal/migrations"); err != nil {
		return fmt.Errorf("failed to run up migrations: %w", err)
	}
	return nil
}

func (a *App) getAuthService(ctx context.Context) (*service.AuthService, error) {
	if a.auth == nil {
		repo, err := a.getRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}
		a.auth = service.New(repo, a.cfg, a.logger)
	}
	return a.auth, nil
}

func (a *App) getAuthServer(ctx context.Context) (*server.Server, error) {
	if a.authServer == nil {
		service, err := a.getAuthService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get service: %w", err)
		}
		a.authServer = server.New(service, a.logger)
	}
	return a.authServer, nil
}

func getGRPCServer(srv *server.Server) *grpc.Server {
	grpcSrv := grpc.NewServer()
	authpb.RegisterAuthServer(grpcSrv, srv)
	return grpcSrv
}

// Close корректно завершает работу приложения
func (a *App) Close() error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}
	return nil
}
