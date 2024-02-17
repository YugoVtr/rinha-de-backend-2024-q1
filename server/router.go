package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
	"github.com/yugovtr/rinha-de-backend-2024-q1/server/log"
)

type App struct {
	Cliente Cliente
}

func Serve(c Cliente) *mux.Router {
	app := &App{Cliente: c}
	router := mux.NewRouter()
	clienteRouter := router.PathPrefix("/clientes/{id:[0-9]+}").Subrouter()
	clienteRouter.HandleFunc("/transacoes", app.Transacoes).Methods("POST")
	clienteRouter.HandleFunc("/extrato", app.Extrato).Methods("GET")
	router.Use(log.LoggingMiddleware)
	return router
}

func (app App) Transacoes(response http.ResponseWriter, request *http.Request) {
	var transacao entity.Transacao
	vars := mux.Vars(request)
	transacao.ClienteID, _ = strconv.ParseInt(vars["id"], 10, 64)
	if !app.Cliente.Existe(transacao.ClienteID) {
		http.Error(response, "cliente não encontrado", http.StatusNotFound)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&transacao)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	if !transacao.Validar() {
		http.Error(response, "dados inválidos", http.StatusBadRequest)
		return
	}
	conta, err := app.Cliente.Transacao(transacao)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(map[string]any{"limite": conta.Limite, "saldo": conta.Saldo}); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app App) Extrato(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	clienteID, _ := strconv.ParseInt(vars["id"], 10, 64)
	if !app.Cliente.Existe(clienteID) {
		http.Error(response, "cliente não encontrado", http.StatusNotFound)
		return
	}
	extrato, err := app.Cliente.Extrato(clienteID)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(extrato); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}
