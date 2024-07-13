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

func run(ctx context.Context) (*api.Api, error) {
	err := config.Load(*cfg)
	if err != nil {
		slog.Error("Reading configuration", "error", err)
		return nil, err
	}
	level, ok := config.MapLevel[strings.ToUpper(config.Confs.Logger.Level)]
	if !ok {
		level = slog.LevelError
	}
	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: config.Confs.Logger.AddSource,
		Level:     level,
	}))
	slog.SetDefault(l)

	server, err := api.New()
	if err != nil {
		slog.Error("Bootin api", "error", err)
		return nil, err
	}

	go func() {
		if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Starting server", "errror", err)
			os.Exit(1)
		}
	}()

	return server, nil
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
			slog.Info("Received SIGHUP, reloading configuration...")

			err = server.Shutdown(ctx)
			if err != nil {
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
