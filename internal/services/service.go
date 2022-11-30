package services

import (
	"context"
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

const (
	AvailableSymbolInID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	IDLength            = 7
)

func (r *Reader) ReadQueue() {
	//start loop reading request channel for every client
	go func() {
		for {
			for _, ch := range r.ReadQ.UserQ {
				go func(ch chan model.Request) {
					for req := range ch {
						money, err := r.WithdrawalMoney(context.Background(), req)
						if err != nil {
							r.Resp.RespQ <- &model.Response{Result: "", Error: model.ResultErr}
							r.logger.Sugar().Errorf("transaction ERROR: %s", err.Error())
						} else {
							r.Resp.RespQ <- money
						}
					}
				}(ch)
			}
		}
	}()
}

func (r *Reader) WithdrawalMoney(ctx context.Context, request model.Request) (*model.Response, error) {
	//get balance from db
	balance, err := r.repo.GetBalance(ctx, request.ClientID)
	if err != nil {
		return nil, errors.Wrap(err, "get balance from db")
	}

	if balance < 0 || balance-request.Amount < 0 {
		return nil, errors.New("not enough money")
	}

	//decrease balance from client for his request
	err = r.repo.BalanceDecrease(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "decrease balance")
	}

	//end transaction
	err = r.repo.TransactionEnd(ctx, request.TransactionID)
	if err != nil {
		return nil, errors.Wrap(err, "transaction end")
	}

	r.logger.Info("make new transaction",
		zap.Any("transaction_id", request.TransactionID))

	return &model.Response{Result: model.ResultOK, Error: ""}, nil
}

func (r *Reader) StartTransaction(ctx context.Context, request model.Request) (string, error) {
	balance, err := r.repo.GetBalance(ctx, request.ClientID)
	if err != nil {
		return "", errors.Wrap(err, "get balance from db")
	}

	transactionID := GetTransactionID()
	trans := &model.Transaction{
		TransactionID:            transactionID,
		ClientID:                 request.ClientID,
		BalanceBeforeTransaction: balance,
		WithdrawalAmount:         request.Amount,
		TransactionEnd:           false,
	}

	err = r.repo.TransactionStart(ctx, trans)
	if err != nil {
		return "", errors.Wrap(err, "start transaction")
	}

	return transactionID, nil
}

func (r *Reader) RecoveryUncompletedTransactions() error {
	ctx := context.Background()
	transactions, err := r.repo.GetUncompletedTransactions(ctx)
	if err != nil {
		return errors.Wrap(err, "get uncompleted transaction")
	}

	for _, val := range transactions {
		balance, err := r.repo.GetBalance(ctx, val.ClientID)
		if err != nil {
			return errors.Wrap(err, "get balance for recovery")
		}

		if balance == val.BalanceBeforeTransaction {
			req := model.Request{
				ClientID:      val.ClientID,
				Amount:        val.WithdrawalAmount,
				TransactionID: val.TransactionID,
			}
			_, err := r.WithdrawalMoney(ctx, req)
			if err != nil {
				return errors.Wrap(err, "recovery withdrawal")
			}
		}
	}

	return nil
}

func GetTransactionID() string {
	rand.Seed(time.Now().UnixNano())
	var key string

	rs := []rune(AvailableSymbolInID)
	lenOfArray := len(rs)

	for i := 0; i < IDLength; i++ {
		key += string(rs[rand.Intn(lenOfArray)])
	}
	return key
}
