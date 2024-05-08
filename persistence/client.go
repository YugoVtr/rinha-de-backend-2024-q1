package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/yugovtr/rinha-de-backend-2024-q1/entity"
	"github.com/yugovtr/rinha-de-backend-2024-q1/tracer"
)

var cache = map[int64]bool{}

type Cliente struct {
	*sql.DB
}

var trace = tracer.New("persistence")

func NewCliente(db *sql.DB) *Cliente {
	return &Cliente{db}
}

func (c *Cliente) StartTrace(ctx context.Context, name string, opts ...tracer.SpanOpts) func(err error) {
	_, span := trace.Start(ctx, name, tracer.GetSpanOpts(opts)...)
	return func(err error) {
		if err != nil {
			span.SetStatus(tracer.SetSpanError(err))
		}
		span.End()
	}
}

func (c *Cliente) Existe(ctx context.Context, clienteID int64) (exists bool) {
	var err error
	callback := c.StartTrace(ctx, "existe", tracer.SpanOpts{"cliente_id": clienteID})
	defer func() { callback(err) }()

	if v, ok := cache[clienteID]; ok {
		return v
	}
	query := `SELECT EXISTS(SELECT 1 FROM contas WHERE cliente_id=$1)`
	err = c.DB.QueryRow(query, clienteID).Scan(&exists)
	cache[clienteID] = exists
	return exists
}

func (c *Cliente) Transacao(ctx context.Context, transacao entity.Transacao) (conta entity.Conta, err error) {
	callback := c.StartTrace(ctx, "transacao", tracer.SpanOpts{
		"cliente_id": transacao.ClienteID,
		"valor":      transacao.Valor,
		"tipo":       transacao.Tipo,
		"descricao":  transacao.Descricao,
	})
	defer func() { callback(err) }()

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

func (c *Cliente) Extrato(ctx context.Context, clienteID int64) (extrato entity.Extrato, err error) {
	callback := c.StartTrace(ctx, "extrato", tracer.SpanOpts{"cliente_id": clienteID})
	defer func() { callback(err) }()

	extrato = entity.Extrato{UltimasTransacoes: []entity.Transacao{}}
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
