package model

import (
	"strings"
	"time"
)

// Customer - Struct that represents a customer
type Customer struct {
	ID           uint64    `json:"id,omitempty"`
	Name         string    `json:"nome,omitempty"`
	AccountLimit int64     `json:"limite,omitempty"`
	Balance      int64     `json:"saldo,omitempty"`
	DateTime     time.Time `json:"datahora,omitempty"`
}

// Statement - Struct that represents a statement
type Statement struct {
	Balance          Balance       `json:"saldo,omitempty"`
	LastTransactions []Transaction `json:"ultimas_transacoes,"`
}

// Saldo - Struct that represents a balance
type Balance struct {
	Total         int64     `json:"total,"`
	DateStatement time.Time `json:"data_extrato,omitempty"`
	AccountLimit  int64     `json:"limite,"`
}

// Transacao - Struct that represents a transaction
type Transaction struct {
	Value       int64     `json:"valor,omitempty"`
	Type        string    `json:"tipo,omitempty"`
	Description string    `json:"descricao,omitempty"`
	DateTime    time.Time `json:"realizada_em,omitempty"`
}

// Validate - Function that validates the field's transaction
func (txBody *Transaction) Validate() string {

	if txBody.Value <= 0 {
		return "valor deve um número inteiro positivo que representa centavos (não vamos trabalhar com frações de centavos). Por exemplo, R$ 10 são 1000 centavos"
	}

	if txBody.Type != "c" && txBody.Type != "d" {
		return "tipo deve ser apenas c para crédito ou d para débito"
	}

	var strLen int = len(strings.TrimSpace(txBody.Description))
	if strLen <= 0 || strLen > 10 {
		return "descricao deve ser uma string de 1 a 10 caractéres"
	}

	return ""
}
