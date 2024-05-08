package server

import (
	"context"

	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
)

type Cliente interface {
	Existe(ctx context.Context, id int64) bool
	Transacao(context.Context, entity.Transacao) (entity.Conta, error)
	Extrato(ctx context.Context, id int64) (entity.Extrato, error)
}
