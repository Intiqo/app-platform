//go:build wireinject

package dependency

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Intiqo/app-platform/internal/database"
	"github.com/Intiqo/app-platform/internal/http/api"
	"github.com/Intiqo/app-platform/internal/http/handler"
	aAws "github.com/Intiqo/app-platform/internal/pkg/cloud/aws"
	"github.com/Intiqo/app-platform/internal/pkg/config"
	"github.com/Intiqo/app-platform/internal/pkg/secrets"
	"github.com/Intiqo/app-platform/internal/repository"
	"github.com/Intiqo/app-platform/internal/service"
)

func NewAWSConfig(profile string) (aws.Config, error) {
	wire.Build(aAws.NewAWSConfig)
	return aws.Config{}, nil
}

// NewConfig returns a new AppConfig
func NewConfig(awsCfg aws.Config, options config.Options) (config.AppConfig, error) {
	wire.Build(
		secrets.NewAWSSecretsManager,
		config.NewConfig,
	)

	return config.AppConfig{}, nil
}

// NewDatabase returns a new database connection pool
func NewDatabase(cfg config.AppConfig) (*pgxpool.Pool, error) {
	wire.Build(
		database.NewDB,
	)

	return &pgxpool.Pool{}, nil
}

// NewAppApi returns a new AppApi
func NewAppApi(cfg config.AppConfig, awsCfg aws.Config, db *pgxpool.Pool) (*api.AppApi, error) {
	// Build the dependency graph
	wire.Build(
		repository.NewTransactioner,
		repository.NewSettingRepository,

		service.NewSettingService,

		handler.NewSettingHandler,

		api.NewAppApi,
	)

	return &api.AppApi{}, nil
}
