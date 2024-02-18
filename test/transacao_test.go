//go:build integration
// +build integration

package test

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
)

// Conta is used to deserialize the service response
type Conta struct {
	entity.Conta
	Saldo int64 `json:"saldo"`
}

// DoTransacao performs a transacao and returns the conta
func DoTransacao(t *testing.T, clienteID, value int64) (conta Conta) {
	t.Helper()

	val, op := uint64(value), "c"
	if value < 0 {
		val, op = uint64(-value), "d"
	}

	transacao := entity.Transacao{
		Valor:     val,
		Tipo:      op,
		Descricao: t.Name(),
	}
	buffer := bytes.Buffer{}
	err := json.NewEncoder(&buffer).Encode(transacao)
	require.NoError(t, err, "json.Encode()")

	response, err := http.Post(fmt.Sprintf("%s/%d/transacoes", serviceURL, clienteID), "application/json", &buffer)
	require.NoError(t, err, "http.Post()")
	t.Cleanup(func() { response.Body.Close() })
	require.Equal(t, http.StatusOK, response.StatusCode, "response.StatusCode")

	err = json.NewDecoder(response.Body).Decode(&conta)
	require.NoError(t, err, "json.Decode()")
	return conta
}

// DoTransacoes performs a list of transacoes
func DoTransacoes(t *testing.T, clienteID int64, transacoes []int64) {
	t.Helper()
	t.Parallel()
	for _, value := range transacoes {
		t.Run(fmt.Sprintf("transacao-(%d)", value), func(t *testing.T) {
			DoTransacao(t, clienteID, value)
		})
	}
}

// GenerateRandomTransacoes generates a list of random transacoes and returns the sum
func GenerateRandomTransacoes(t *testing.T, n int) ([]int64, int64) {
	t.Helper()

	transacoes, sum := make([]int64, n, n), int64(0)
	for i := 0; i < n; i++ {
		var value int64
		binary.Read(rand.Reader, binary.LittleEndian, &value)
		value = value % 10_000
		if value%2 == 0 {
			value = -value
		}
		transacoes[i] = value
		sum += value
	}
	return transacoes, sum
}
