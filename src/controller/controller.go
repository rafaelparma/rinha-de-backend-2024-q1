package controller

import (
	"encoding/json"
	"math"
	"rinha-de-backend-2024-q1/src/database"
	"rinha-de-backend-2024-q1/src/model"
	"rinha-de-backend-2024-q1/src/response"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

type poolDB struct {
	db *database.DB
}

// InjectDB - Function that injects the database
func InjectDB(db *database.DB) *poolDB {
	return &poolDB{db: db}
}

// Statement - Function that returns the customer's transactions history and balance
func (poolDB *poolDB) Statement(ctx *fasthttp.RequestCtx) {

	customer_id, err := strconv.ParseUint(ctx.UserValue("customer_id").(string), 10, 64)
	if err != nil {
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, "erro ao ler o id do cliente")
		return
	}

	tx, err := poolDB.db.DB().BeginTx(ctx, nil)
	if err != nil {
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, "erro na transacao")
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		"SELECT id, name, account_limit, balance, datetime FROM customers WHERE id = $1",
		customer_id)

	var customer model.Customer
	err = row.Scan(&customer.ID, &customer.Name, &customer.AccountLimit, &customer.Balance, &customer.DateTime)
	if err != nil {
		tx.Rollback()
		response.RError(ctx, fasthttp.StatusNotFound, "id do cliente não encontrado")
		return
	}
	rows, err := tx.Query(
		"SELECT value, type, description, datetime FROM transactions WHERE customer_id = $1 ORDER BY id DESC LIMIT 10",
		customer_id)
	if err != nil {
		tx.Rollback()
		response.RError(ctx, fasthttp.StatusNotFound, "erro ao consultar as transações do cliente")
		return
	}

	var transactions []model.Transaction
	for rows.Next() {
		var row model.Transaction
		err = rows.Scan(&row.Value, &row.Type, &row.Description, &row.DateTime)
		if err != nil {
			tx.Rollback()
			response.RError(ctx, fasthttp.StatusNotFound, "erro ao consultar as transações do cliente")
			return
		}

		transactions = append(transactions, row)
	}
	tx.Commit()

	var data = model.Statement{
		Balance: model.Balance{
			Total:         customer.Balance,
			DateStatement: time.Now(),
			AccountLimit:  customer.AccountLimit,
		},
		LastTransactions: transactions,
	}

	response.RJSON(ctx, fasthttp.StatusOK, data)

}

// Transactions - Function that insert the transaction in the customer's account
func (poolDB *poolDB) Transactions(ctx *fasthttp.RequestCtx) {

	bodyRequest := ctx.Request.Body()

	var txBody model.Transaction
	if err := json.Unmarshal(bodyRequest, &txBody); err != nil {
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, "erro ao ler o corpo da requisição")
		return
	}

	if err := txBody.Validate(); err != "" {
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, err)
		return
	}

	customer_id, err := strconv.ParseUint(ctx.UserValue("customer_id").(string), 10, 64)
	if err != nil {
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, "erro ao ler o id do cliente")
		return
	}

	tx, err := poolDB.db.DB().BeginTx(ctx, nil)
	if err != nil {
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, "erro na transacao")
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		"SELECT id, name, account_limit, balance, datetime FROM customers WHERE id = $1 FOR UPDATE",
		customer_id)

	var customer model.Customer
	err = row.Scan(&customer.ID, &customer.Name, &customer.AccountLimit, &customer.Balance, &customer.DateTime)
	if err != nil {
		response.RError(ctx, fasthttp.StatusNotFound, "id do cliente não encontrado")
		return
	}

	newBalance := customer.Balance + txBody.Value
	if txBody.Type == "d" {
		if float64(customer.AccountLimit) < math.Abs(float64(customer.Balance-txBody.Value)) {
			tx.Rollback()
			response.RError(ctx, fasthttp.StatusUnprocessableEntity, "saldo insuficiente")
			return
		}
		newBalance = customer.Balance - txBody.Value
	}

	_, err = tx.Exec(
		"INSERT INTO transactions (customer_id, value, type, description) values ($1, $2, $3, $4)",
		customer.ID, txBody.Value, txBody.Type, txBody.Description)
	if err != nil {
		tx.Rollback()
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, "erro na transacao")
		return
	}

	updatecmd := "UPDATE customers SET balance = (balance + $1) WHERE id = $2"
	if txBody.Type == "d" {
		updatecmd = "UPDATE customers SET balance = (balance - $1) WHERE id = $2"
	}
	_, err = tx.Exec(updatecmd, txBody.Value, customer.ID)
	if err != nil {
		tx.Rollback()
		response.RError(ctx, fasthttp.StatusUnprocessableEntity, "saldo insuficiente")
		return
	}

	tx.Commit()

	var data = struct {
		AccountLimit int64 `json:"limite"`
		Balance      int64 `json:"saldo"`
	}{

		AccountLimit: customer.AccountLimit,
		Balance:      newBalance,
	}

	response.RJSON(ctx, fasthttp.StatusOK, data)

}
