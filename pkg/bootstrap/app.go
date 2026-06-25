package bootstrap

import (
	"context"
	"fmt"

	locationsgrpc "github.com/zchelalo/neuraclinic-location/internal/modules/locations/adapters/grpc"
	locationpg "github.com/zchelalo/neuraclinic-location/internal/modules/locations/adapters/persistence/postgres"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application"
	grpcserver "github.com/zchelalo/neuraclinic-location/internal/server/grpc"
	"go.uber.org/zap"
)

type App struct {
	Server  *grpcserver.Server
	Cleanup func(context.Context) error
}

func InitApp(ctx context.Context, logger *zap.Logger, cfg Config) (*App, error) {
	db, err := NewDB(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db: %w", err)
	}

	repo := locationpg.NewRepository(db)
	locationApp := application.NewService(application.Config{
		DefaultCountryCode: cfg.LocationDefaultCountryCode,
		LimitDefault:       cfg.LocationLimitDefault,
		LimitMax:           cfg.LocationLimitMax,
	}, repo)

	server, err := grpcserver.New(grpcserver.Config{
		Port:            cfg.Port,
		ServiceName:     cfg.ServiceName,
		TLSCertFilePath: cfg.GRPCTLSCertPath,
		TLSKeyFilePath:  cfg.GRPCTLSKeyPath,
	}, logger, grpcserver.Services{
		Location: locationsgrpc.NewLocationService(locationApp),
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot create grpc server: %w", err)
	}

	return &App{
		Server: server,
		Cleanup: func(context.Context) error {
			server.GracefulStop()
			db.Close()
			return nil
		},
	}, nil
}
