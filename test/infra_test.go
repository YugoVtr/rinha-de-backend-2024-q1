//go:build integration
// +build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
	"github.com/yugovtr/rinha-de-backend-2024-q1/persistence"
)

func StartStack(t *testing.T) {
	t.Helper()

	compose, err := tc.NewDockerCompose("../compose.yml")
	require.NoError(t, err, "NewComposeDB()")
	t.Cleanup(func() {
		err := compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)
		require.NoError(t, err, "compose.Down()")
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")
	time.Sleep(5 * time.Second)
}

func ResetDatabase(t *testing.T) {
	t.Helper()

	sqlDB, err := persistence.ConnectToPostgres(t.Logf)
	require.NoError(t, err, "ConnectToPostgres()")

	_, err = sqlDB.Exec("TRUNCATE TABLE transacoes")
	require.NoError(t, err, "Failed to truncate table transacoes")

	seeds := []entity.Conta{
		{ClienteID: 1, Saldo: 0, Limite: 100000},
		{ClienteID: 2, Saldo: 0, Limite: 80000},
		{ClienteID: 3, Saldo: 0, Limite: 1000000},
		{ClienteID: 4, Saldo: 0, Limite: 10000000},
		{ClienteID: 5, Saldo: 0, Limite: 500000},
	}
	for _, seed := range seeds {
		_, err = sqlDB.Exec("UPDATE contas SET total=$1, limite=$2 WHERE cliente_id=$3", seed.Saldo, seed.Limite, seed.ClienteID)
		require.NoError(t, err, "Failed to insert seed")
	}
}
