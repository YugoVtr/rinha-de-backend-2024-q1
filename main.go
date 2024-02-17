package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/yugovtr/rinha-de-backend-2024-q1/persistence"
	"github.com/yugovtr/rinha-de-backend-2024-q1/server"
)

func main() {
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
