package main

import (
	"errors"
	"github.com/vaberof/notification-grpc/internal/domain/notification"
	"github.com/vaberof/notification-grpc/pkg/config"
	"github.com/vaberof/notification-grpc/pkg/grpc/grpcserver"
	"os"
)

type AppConfig struct {
	Server     grpcserver.ServerConfig
	SmtpConfig notification.SmtpConfig
}

func mustGetAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var serverConfig grpcserver.ServerConfig
	err := config.ParseConfig(provider, "app.grpc.server", &serverConfig)
	if err != nil {
		return nil, err
	}

	var smtpConfig notification.SmtpConfig
	smtpConfig.SenderEmail = os.Getenv("SMTP_SENDER")
	smtpConfig.Password = os.Getenv("SMTP_PASSWORD")

	appConfig := AppConfig{
		Server:     serverConfig,
		SmtpConfig: smtpConfig,
	}

	return &appConfig, nil
}
