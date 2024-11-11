package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/Intiqo/app-platform/internal/pkg/secrets"
)

const SourceKey = "CONFIG_SOURCE"

const SourceEnv = "ENVIRONMENT"
const SourceAWSSecretsManager = "AWS_SECRETS_MANAGER"

const AwsProfileKey = "AWS_PROFILE"
const AwsConfigSecretsNameKey = "AWS_CONFIG_SECRETS_NAME"

type Options struct {
	ConfigSource  string
	ConfigFile    string
	AwsProfile    string
	AWSecretsName string
}

type AppConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort int    `mapstructure:"APP_PORT"`

	AuthSecret       string `mapstructure:"AUTH_SECRET"`
	AuthExpiryPeriod int    `mapstructure:"AUTH_EXPIRY_PERIOD"`

	DatabaseHost     string `mapstructure:"DB_HOST"`
	DatabasePort     string `mapstructure:"DB_PORT"`
	DatabaseUsername string `mapstructure:"DB_USERNAME"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	DatabaseName     string `mapstructure:"DB_DATABASE_NAME"`

	RequestBodySizeLimit string `mapstructure:"REQUEST_BODY_SIZE_LIMIT"`

	SwaggerHostUrl    string `mapstructure:"SWAGGER_HOST_URL"`
	SwaggerHostScheme string `mapstructure:"SWAGGER_HOST_SCHEME"`
	SwaggerUsername   string `mapstructure:"SWAGGER_USERNAME"`
	SwaggerPassword   string `mapstructure:"SWAGGER_PASSWORD"`
}

type configManager struct {
	sm secrets.Manager
}

// NewConfig returns a new AppConfig
func NewConfig(opts Options, sm secrets.Manager) (AppConfig, error) {
	cm := &configManager{
		sm: sm,
	}
	switch opts.ConfigSource {
	case SourceEnv:
		return cm.newFromEnvironment(opts)
	case SourceAWSSecretsManager:
		return cm.newFromAWSSecretsManager(opts)
	}
	return AppConfig{}, errors.New("invalid config type")
}

func (m *configManager) newFromEnvironment(opts Options) (cfg AppConfig, err error) {
	viper.SetConfigFile(opts.ConfigFile)
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %v", err)
	}

	cfg = AppConfig{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to load configuration: %v", err)
	}

	return cfg, nil
}

func (m *configManager) newFromAWSSecretsManager(opts Options) (cfg AppConfig, err error) {
	secret, err := m.sm.GetSecret(opts.AWSecretsName)
	if err != nil {
		return cfg, fmt.Errorf("failed to load configuration from AWS Secrets Manager: %v", err)
	}

	viper.SetConfigType("json")
	err = viper.ReadConfig(strings.NewReader(secret))
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %v", err)
	}

	cfg = AppConfig{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to load configuration: %v", err)
	}

	return cfg, nil
}
