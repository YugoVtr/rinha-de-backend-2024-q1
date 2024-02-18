//go:build integration
// +build integration

package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
)

func GetExtrato(t *testing.T, clienteID int64) (extrato entity.Extrato) {
	t.Helper()

	response, err := http.Get(fmt.Sprintf("%s/%d/extrato", serviceURL, clienteID))
	require.NoError(t, err, "http.Get()")
	t.Cleanup(func() { response.Body.Close() })
	require.Equal(t, http.StatusOK, response.StatusCode, "response.StatusCode")

	err = json.NewDecoder(response.Body).Decode(&extrato)
	require.NoError(t, err, "json.Decode()")
	return extrato
}
