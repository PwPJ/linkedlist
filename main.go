package main

import (
	"context"
	"errors"
	"flag"
	"linkedlist/api"
	"linkedlist/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var cfg = flag.String("config", "config/config.yaml", "yaml config path")

type ConfigChangeType string

const (
	serverChange ConfigChangeType = "server"
	loggerChange ConfigChangeType = "logger"
	noChange     ConfigChangeType = "no change"
)

func run(ctx context.Context) (*api.Api, error) {
	err := config.Load(*cfg)
	if err != nil {
		slog.Error("Reading configuration", "error", err)
		return nil, err
	}

	ConfigureLoggerLevel()

	server, err := api.New()
	if err != nil {
		slog.Error("Booting api", "error", err)
		return nil, err
	}

	go func() {
		if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Starting server", "error", err)
			os.Exit(1)
		}
	}()

	return server, nil
}

func ConfigureLoggerLevel() {
	level, ok := config.MapLevel[strings.ToUpper(config.Confs.Logger.Level)]
	if !ok {
		level = slog.LevelError
	}

	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: config.Confs.Logger.AddSource,
		Level:     level,
	}))

	slog.SetDefault(l)
}

func configChanged(oldConfig *config.Config) ConfigChangeType {
	if oldConfig.Server.Port != config.Confs.Server.Port {
		return serverChange
	}

	if oldConfig.Logger.AddSource != config.Confs.Logger.AddSource ||
		oldConfig.Logger.Level != config.Confs.Logger.Level {
		return loggerChange
	}

	return noChange
}

func main() {
	flag.Parse()

	ctx := context.Background()

	server, err := run(ctx)
	if err != nil {
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGHUP:
			slog.Info("Received SIGHUP, checking configuration...")

			oldConfig := config.Confs

			if err := config.Load(*cfg); err != nil {
				slog.Error("Reading new configuration", "error", err)
				continue
			}

			configChangeType := configChanged(&oldConfig)

			if configChangeType == noChange {
				slog.Info("Configuration unchanged, no actions required.")
				continue
			}

			if configChangeType == loggerChange {
				ConfigureLoggerLevel()
				slog.Info("Logger configuration updated.")
				continue
			}

			slog.Info("Configuration has changed, reloading server...")

			if err = server.Shutdown(ctx); err != nil {
				slog.Error("could not stop server", "error", err)
				continue
			}

			server, err = run(ctx)
			if err != nil {
				os.Exit(1)
			}

		case syscall.SIGINT, syscall.SIGTERM:
			slog.Info("Received SIGINT/SIGTERM, shutting down...")
			if err := server.Shutdown(ctx); err != nil {
				slog.Error("could not gracefully shut down server", "error", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}
