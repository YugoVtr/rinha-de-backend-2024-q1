package entity

type Extrato struct {
	Saldo             Conta       `json:"saldo"`
	UltimasTransacoes []Transacao `json:"ultimas_transacoes"`
}
