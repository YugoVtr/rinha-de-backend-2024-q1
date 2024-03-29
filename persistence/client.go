package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
)

var cache = map[int64]bool{}

type Cliente struct {
	*sql.DB
}

func NewCliente(db *sql.DB) *Cliente {
	return &Cliente{db}
}

func (c *Cliente) Existe(clienteID int64) (exists bool) {
	if v, ok := cache[clienteID]; ok {
		return v
	}
	query := `SELECT EXISTS(SELECT 1 FROM contas WHERE cliente_id=$1)`
	_ = c.DB.QueryRow(query, clienteID).Scan(&exists)
	cache[clienteID] = exists
	return exists
}

func (c *Cliente) Transacao(transacao entity.Transacao) (conta entity.Conta, err error) {
	trasaction, err := c.DB.Begin()
	if err != nil {
		return entity.Conta{}, fmt.Errorf("erro ao iniciar transação %w", err)
	}
	defer func() {
		if err != nil {
			if txErr := trasaction.Rollback(); txErr != nil {
				err = errors.Join(err, txErr)
			}
			return
		}
		if commitErr := trasaction.Commit(); commitErr != nil {
			err = errors.Join(err, commitErr)
			if txErr := trasaction.Rollback(); txErr != nil {
				err = errors.Join(err, txErr)
			}
		}
	}()

	query := `SELECT cliente_id, total, limite FROM contas WHERE cliente_id=$1`
	if err := trasaction.QueryRow(query, transacao.ClienteID).Scan(&conta.ClienteID, &conta.Saldo, &conta.Limite); err != nil {
		return conta, fmt.Errorf("conta não encontrada %w", err)
	}

	novaConta, err := conta.Exec(transacao)
	if err != nil {
		return novaConta, fmt.Errorf("erro ao executar transação %w", err)
	}
	query = `UPDATE contas SET total=$1 WHERE cliente_id=$2`
	if _, err = trasaction.Exec(query, novaConta.Saldo, novaConta.ClienteID); err != nil {
		return novaConta, fmt.Errorf("erro ao atualizar saldo %w", err)
	}

	query = `INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizado_em) VALUES ($1, $2, $3, $4, $5)`
	if _, err = trasaction.Exec(query, transacao.ClienteID, transacao.Valor, transacao.Tipo, transacao.Descricao, time.Now()); err != nil {
		return novaConta, fmt.Errorf("erro ao inserir transação %w", err)
	}
	return novaConta, nil
}

func (c *Cliente) Extrato(clienteID int64) (entity.Extrato, error) {
	extrato := entity.Extrato{UltimasTransacoes: []entity.Transacao{}}
	query := `SELECT total, limite FROM contas WHERE cliente_id=$1`
	if err := c.DB.QueryRow(query, clienteID).Scan(&extrato.Saldo.Saldo, &extrato.Saldo.Limite); err != nil {
		return extrato, fmt.Errorf("conta não encontrada %w", err)
	}
	extrato.Saldo.DataExtrato = time.Now().Local().Format(time.RFC3339Nano)

	query = `SELECT valor, tipo, descricao, realizado_em FROM transacoes WHERE cliente_id=$1 ORDER BY realizado_em DESC LIMIT 10`
	rows, err := c.DB.Query(query, clienteID)
	if err != nil {
		return extrato, fmt.Errorf("erro na consulta ao buscar transações %w", err)
	}
	if rows.Err() != nil {
		return extrato, fmt.Errorf("erro no resultado ao buscar transações %w", rows.Err())
	}
	defer rows.Close()

	for rows.Next() {
		transacao := entity.Transacao{}
		if err := rows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadoEm); err != nil {
			return extrato, fmt.Errorf("erro ao escanear transação %w", err)
		}
		extrato.UltimasTransacoes = append(extrato.UltimasTransacoes, transacao)
	}

	return extrato, nil
}
