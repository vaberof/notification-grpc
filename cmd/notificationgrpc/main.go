package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vaberof/notification-grpc/internal/app/entrypoint/grpc/notification"
	notificationservice "github.com/vaberof/notification-grpc/internal/domain/notification"
	"github.com/vaberof/notification-grpc/pkg/grpc/grpcserver"
	"github.com/vaberof/notification-grpc/pkg/logging/logs"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "List of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()

	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	logger := logs.New(os.Stdout, nil)

	appConfig := mustGetAppConfig(*appConfigPaths)

	fmt.Printf("%+v\n", appConfig)

	grpcServer := grpcserver.New(&appConfig.Server, logger)

	notificationService := notificationservice.NewNotificationService(&appConfig.SmtpConfig, logger)

	notification.Register(grpcServer.Server, notificationService)

	grpcServerErrorCh := grpcServer.StartAsync()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case signalValue := <-quitCh:
		logger.GetLogger().Info("stopping application", slog.String("signal", signalValue.String()))

		grpcServer.Shutdown()
	case err := <-grpcServerErrorCh:
		logger.GetLogger().Info("stopping application", slog.String("gRPC server error", err.Error()))

		grpcServer.Shutdown()
	}
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
