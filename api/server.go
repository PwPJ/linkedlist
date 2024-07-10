package api

import (
	"context"
	"fmt"
	v1 "linkedlist/api/v1"
	v2 "linkedlist/api/v2"

	// v2 "linkedlist/api/v2"
	// "linkedlist/config"
	"log/slog"
	"net/http"
)

type Api struct {
	Mux    *http.ServeMux
	Server *http.Server
}

func New() (*Api, error) {
	v1 := v1.V1()
	v2, err := v2.V2()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/v1/", http.StripPrefix("/v1", v1))
	mux.Handle("/v2/", http.StripPrefix("/v2", v2))

	return &Api{
		Mux: mux,
	}, nil
}

func (a *Api) Shutdown(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}

func (a *Api) Start() error {
	port := 8080
	a.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Mux,
	}
	slog.Info("Running", "port", port)
	err := a.Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
