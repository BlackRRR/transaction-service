package repository

import (
	"context"
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type Implementation interface {
	GetBalance(ctx context.Context, clientID int64) (int64, error)
	BalanceDecrease(ctx context.Context, request model.TransactionReq) error

	GetUncompletedTransactions(ctx context.Context, status string) ([]*model.Transaction, error)
	TransactionStart(ctx context.Context, transaction *model.Transaction) error
	ChangeTransactionStatus(ctx context.Context, transactionID string, status string) error
	GetAllTransactions(ctx context.Context, clientID int64) ([]*model.Transaction, error)
}
type Repository struct {
	db *pgxpool.Pool
}

func NewRepo(ctx context.Context, connect *pgxpool.Pool) (*Repository, error) {
	repo := &Repository{
		db: connect,
	}

	rows, err := repo.db.Query(ctx, `
CREATE TABLE IF NOT EXISTS transactions(
	transaction_id text UNIQUE, 
	client_id bigint,
	balance_before bigint,
	withdrawal_amount bigint,
	transaction_status text);`)
	if err != nil {
		return nil, errors.Wrap(err, "create transaction table")
	}
	defer rows.Close()

	rows, err = repo.db.Query(ctx, `
CREATE TABLE IF NOT EXISTS client(
	client_id bigint,
	balance bigint);`)
	if err != nil {
		return nil, errors.Wrap(err, "create transaction table")
	}

	return repo, nil
}

func (r *Repository) TransactionStart(ctx context.Context, transaction *model.Transaction) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO transactions (
                         transaction_id, 
                         client_id,
                         balance_before,
                         withdrawal_amount,
                         transaction_status) 
VALUES ($1,$2,$3,$4,$5)`,
		transaction.TransactionID,
		transaction.ClientID,
		transaction.BalanceBeforeTransaction,
		transaction.WithdrawalAmount,
		transaction.TransactionStatus)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *Repository) ChangeTransactionStatus(ctx context.Context, transactionID string, status string) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to begin end transaction")
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE transactions SET transaction_status = $1 WHERE transaction_id = $2`,
		status,
		transactionID)
	if err != nil {
		return errors.Wrap(err, "failed to end transaction")
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *Repository) BalanceDecrease(ctx context.Context, request model.TransactionReq) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction balance decrease")
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE client SET balance = balance - $1 WHERE client_id = $2`,
		request.Amount,
		request.ClientID)
	if err != nil {
		return errors.Wrap(err, "failed to update balance to client")
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *Repository) GetBalance(ctx context.Context, clientID int64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "failed to begin transaction get balance")

	}

	defer tx.Rollback(ctx)

	var balance int64
	err = tx.QueryRow(ctx, `SELECT balance FROM client WHERE client_id = $1`,
		clientID).Scan(&balance)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return 0, errors.New("client does not exists")
		}
		return 0, errors.Wrap(err, "failed to get balance from client")
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return balance, nil
}

func (r *Repository) GetUncompletedTransactions(ctx context.Context, status string) ([]*model.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction get uncompleted")
	}

	defer tx.Rollback(ctx)

	rows, err := r.db.Query(ctx, `SELECT
    	transaction_id,
       	client_id,
       	balance_before,
       	withdrawal_amount,
       	transaction_status FROM transactions
                         WHERE transaction_status = $1`, status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get uncompleted transactions")
	}
	defer rows.Close()

	transactions, err := readRows(rows)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return transactions, nil
}

func (r *Repository) GetAllTransactions(ctx context.Context, clientID int64) ([]*model.Transaction, error) {
	rows, err := r.db.Query(ctx, `SELECT
       transaction_id,
       client_id,
       balance_before,
       withdrawal_amount,
       transaction_status FROM transactions 
                          WHERE client_id = $1`, clientID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all transactions")
	}

	transactions, err := readRows(rows)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func readRows(rows pgx.Rows) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	for rows.Next() {
		transaction := &model.Transaction{}
		err := rows.Scan(
			&transaction.TransactionID,
			&transaction.ClientID,
			&transaction.BalanceBeforeTransaction,
			&transaction.WithdrawalAmount,
			&transaction.TransactionStatus)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan transactions")
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
