package main

import (
	"linkedlist/api"
	"log/slog"
	"os"
)

func main() {
	server, err := api.New()
	if err != nil {
		slog.Error("Bootin api", "error", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		slog.Error("Starting server", "errror", err)
		os.Exit(1)
	}

}
