package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"

	"github.com/Intiqo/app-platform/internal/dependency"
	"github.com/Intiqo/app-platform/internal/http/swagger"
	"github.com/Intiqo/app-platform/internal/pkg/config"
	"github.com/Intiqo/app-platform/internal/version"
)

func main() {
	// Print the current api version
	version.PrintInfo()

	// Get the configuration options
	cfgOptions := getConfigOptions()

	awsCfg, err := dependency.NewAWSConfig(cfgOptions.AwsProfile)
	if err != nil {
		log.Fatalf("failed to load aws config: %v", err)
	}

	// Validate the aws session
	validateAwsSession(awsCfg)

	// Load the configuration
	cfg, err := dependency.NewConfig(awsCfg, cfgOptions)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Setup database connection
	db, err := dependency.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to create connection for database: %v", err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// Defer closing the database connection
	defer db.Close()

	// Initialize the dependencies
	api, err := dependency.NewAppApi(
		cfg, awsCfg, db,
	)
	if err != nil {
		log.Fatalf("failed to create app api: %v", err)
	}

	// Set up the echo server
	e := echo.New()
	e.HideBanner = true

	// Set up the middleware
	api.SetupMiddleware(e)

	// Set up the swagger documentation
	swagger.SetupSwagger(cfg, e)

	// Set up the routes
	api.SetupRoutes(e)

	// Start the server in a goroutine to handle graceful shutdown
	go func() {
		e.Logger.Info(e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.AppPort)))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}

func getConfigOptions() config.Options {
	// Load the configuration
	cfgSource := os.Getenv(config.SourceKey)
	if cfgSource == "" {
		cfgSource = config.SourceEnv
	}
	cfgOptions := config.Options{
		ConfigSource: cfgSource,
	}
	switch cfgSource {
	case config.SourceEnv:
		cfgOptions.ConfigFile = ".env"
		cfgOptions.AwsProfile = os.Getenv(config.AwsProfileKey)
	case config.SourceAWSSecretsManager:
		cfgOptions.AwsProfile = os.Getenv(config.AwsProfileKey)
		cfgOptions.AWSecretsName = os.Getenv(config.AwsConfigSecretsNameKey)
	}
	return cfgOptions
}

func validateAwsSession(cfg aws.Config) {
	// Create an STS client
	svc := sts.NewFromConfig(cfg)

	// Call GetCallerIdentity to verify the session
	result, err := svc.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	if err != nil {
		log.Fatalf("failed to verify the aws session: %v", err)
	}

	log.Printf("aws session verified. Account: %s, Arn: %s\n", *result.Account, *result.Arn)
}
