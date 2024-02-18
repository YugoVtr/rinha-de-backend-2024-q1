//go:build integration
// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const serviceURL = "http://0.0.0.0:9999/clientes"

func TestResetDB(t *testing.T) {
	ResetDatabase(t)
}

func TestIntegration(t *testing.T) {
	StartStack(t)
	ResetDatabase(t)

	const (
		clienteID          = int64(1)
		NumberOfTransacoes = 100
		MaxTransacoes      = 10
	)

	saldo := GetExtrato(t, clienteID).Saldo.Saldo
	transacoes, soma := GenerateRandomTransacoes(t, NumberOfTransacoes)
	DoTransacoes(t, clienteID, transacoes)

	extrato := GetExtrato(t, clienteID)
	require.Len(t, extrato.UltimasTransacoes, MaxTransacoes)
	require.Equal(t, saldo+soma, extrato.Saldo.Saldo)
}
