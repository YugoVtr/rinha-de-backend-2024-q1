package server

import "github.com/yugovtr/rinha-de-backend-2024-q1/entity"

type Cliente interface {
	Existe(id int64) bool
	Transacao(entity.Transacao) (entity.Conta, error)
	Extrato(id int64) (entity.ExtratoResponse, error)
}
