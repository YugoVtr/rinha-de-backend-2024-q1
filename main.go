package main

import (
	"net/http"
	"time"

	"github.com/yugovtr/rinha-de-backend-2024-q1/server"
)

func main() {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        server.Serve(nil),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = server.ListenAndServe()
}

func Transacoes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
