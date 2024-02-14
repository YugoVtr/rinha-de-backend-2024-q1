package server_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
	"github.com/yugovtr/rinha-de-backend-2024-q1/server"
)

func TestTransacoes(t *testing.T) {
	testCases := []struct {
		id      uint64
		payload string
		status  int
	}{
		{1, `{"valor":1000,"tipo":"c","descricao":"descricao"}`, http.StatusOK},
		{2, `{"valor":1000,"tipo":"d","descricao":"descricao"}`, http.StatusOK},
		{6, `{"valor":1000,"tipo":"c","descricao":"descricao"}`, http.StatusNotFound},
		{1, `{"valor":-10,"tipo":"c","descricao":"descricao"}`, http.StatusBadRequest},
		{1, `{"valor":0,"tipo":"c","descricao":"descricao"}`, http.StatusBadRequest},
		{1, `{"valor":1000,"tipo":"a","descricao":"descricao"}`, http.StatusBadRequest},
		{1, `{"valor":1000,"tipo":"c","descricao":""}`, http.StatusBadRequest},
		{5, `{"valor":1000,"tipo":"c","descricao":"descricao"}`, http.StatusUnprocessableEntity},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprint(testCase.id), func(t *testing.T) {
			url := fmt.Sprintf(Setup(t, "transacoes"), testCase.id)
			payload := strings.NewReader(testCase.payload)
			req, err := http.NewRequestWithContext(context.TODO(), "POST", url, payload)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			t.Cleanup(func() { resp.Body.Close() })
			require.NoError(t, err)

			assert.Equal(t, testCase.status, resp.StatusCode, testCase)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			t.Logf("%s", body)
		})
	}
}

func TestExtrato(t *testing.T) {
	testCases := []struct {
		id     uint64
		status int
	}{
		{1, http.StatusOK},
		{6, http.StatusNotFound},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprint(testCase.id), func(t *testing.T) {
			url := fmt.Sprintf(Setup(t, "extrato"), testCase.id)
			req, err := http.NewRequestWithContext(context.TODO(), "GET", url, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			t.Cleanup(func() { resp.Body.Close() })
			require.NoError(t, err)

			assert.Equal(t, testCase.status, resp.StatusCode, testCase)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			t.Logf("%s", body)
		})
	}
}

type client struct{}

func (client) Existe(id int64) bool { return id > 0 && id < 6 }

func (client) Transacao(t entity.Transacao) (entity.Conta, error) {
	if t.ClienteID == 5 {
		return entity.Conta{}, fmt.Errorf("saldo inconsistente")
	}
	return entity.Conta{}, nil
}

func (c client) Extrato(id int64) (entity.Extrato, error) {
	extrato := entity.Extrato{UltimasTransacoes: []entity.Transacao{}}
	if !c.Existe(id) {
		return extrato, fmt.Errorf("cliente não encontrado")
	}
	return extrato, nil
}

func Setup(t *testing.T, action string) string {
	t.Helper()
	s := httptest.NewServer(server.Serve(&client{}))
	t.Cleanup(func() { s.Close() })
	return s.URL + "/clientes/%d/" + action
}
