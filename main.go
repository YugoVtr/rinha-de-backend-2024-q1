package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/yugovtr/rinha-de-backend-2024-q1/persistence"
	"github.com/yugovtr/rinha-de-backend-2024-q1/server"
)

func main() {
	sqlDB, err := persistence.ConnectToPostgres()
	if err != nil {
		slog.Error("connect to postgres", "err", err)
		return
	}
	persistence.NewCliente(sqlDB)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        server.Serve(persistence.NewCliente(sqlDB)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = server.ListenAndServe()
}
