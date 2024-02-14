package entity

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

type Conta struct {
	Limite uint64 `json:"limite"`
	Saldo  int64  `json:"saldo"`
}
