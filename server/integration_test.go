//go:build integration
// +build integration

package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestIntegration(t *testing.T) {
	StartStack(t)
	testCases := []struct {
		id     string
		status int
	}{
		{"1", http.StatusOK},
		{"6", http.StatusNotFound},
	}
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			time.Sleep(5 * time.Second)
			response, err := http.Get("http://0.0.0.0:9999/clientes/" + testCase.id + "/extrato")
			require.NoError(t, err, "http.Get()")
			t.Cleanup(func() { response.Body.Close() })
			require.Equal(t, testCase.status, response.StatusCode, "response.StatusCode")
		})
	}
}

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
}
