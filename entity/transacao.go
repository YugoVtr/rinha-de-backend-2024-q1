package entity

import (
	"encoding/json"
	"errors"
)

var TiposValidosTransacao = map[string]bool{
	"c": true,
	"d": true,
}

type Transacao struct {
	ClienteID   int64  `json:"cliente_id,omitempty"`
	Valor       uint64 `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadoEm string `json:"realizado_em,omitempty"`
}

func (t Transacao) Validar() bool {
	return t.Descricao != "" && t.Valor > 0 && TiposValidosTransacao[t.Tipo]
}

func (t Transacao) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type Conta struct {
	ClienteID   int64  `json:"cliente_id,omitempty"`
	Limite      uint64 `json:"limite"`
	Saldo       int64  `json:"total"`
	DataExtrato string `json:"data_extrato,omitempty"`
}

func (c Conta) Exec(transacao Transacao) (Conta, error) {
	if transacao.Tipo == "d" {
		if c.Saldo-int64(transacao.Valor) < int64(c.Limite)*(-1) {
			return c, errors.New("limite insuficiente")
		}
		c.Saldo -= int64(transacao.Valor)
		return c, nil
	}

	c.Saldo += int64(transacao.Valor)
	return c, nil
}
