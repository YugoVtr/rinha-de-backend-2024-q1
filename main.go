package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yugovtr/rinha-de-backend-2024-q1/persistence"
	"github.com/yugovtr/rinha-de-backend-2024-q1/server"
	"github.com/yugovtr/rinha-de-backend-2024-q1/tracer"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := tracer.InitProvider(os.Getenv("COLLECTOR_URL"), os.Getenv("SERVICE_NAME"))
	if err != nil {
		slog.Error("init tracer provider", "err", err)
		return
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			slog.Error("failed to shutdown TracerProvider", "err", err)
			return
		}
	}()

	sqlDB, err := persistence.ConnectToPostgres(slog.Info)
	if err != nil {
		slog.Error("connect to postgres", "err", err)
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:           net.JoinHostPort("", port),
		Handler:        server.Serve(persistence.NewCliente(sqlDB)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	slog.Info("server started", "port", port)
	defer slog.Info("server finished")
	_ = server.ListenAndServe()
}
