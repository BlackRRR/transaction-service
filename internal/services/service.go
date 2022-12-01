package services

import (
	"context"
	"fmt"
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
)

const (
	AvailableSymbolInID = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	IDLength            = 7
	TransactionStart    = "TRANSACTION IN PROCESS"
	TransactionEnd      = "TRANSACTION PASSED SUCCESSFUL"
	TransactionERROR    = "TRANSACTION ERROR"
)

func (s *TransactionService) ReadQueue() {
	//start loop reading request channel for every client
	go func() {
		for {
			for _, ch := range s.ReadQ.UserQ {
				go func(ch chan model.TransactionReq) {
					for req := range ch {
						money, err := s.WithdrawalMoney(context.Background(), req)
						if err != nil {
							if strings.Contains(err.Error(), "client does not exist") {
								s.Resp.RespQ <- &model.Response{Result: "", Error: "CLIENT DOES NOT EXISTS"}
								continue
							}

							if strings.Contains(err.Error(), "not enough money") {
								s.Resp.RespQ <- &model.Response{Result: "", Error: err.Error()}
								continue
							}

							s.Resp.RespQ <- &model.Response{Result: "", Error: model.ResultErr}
							s.logger.Sugar().Errorf("transaction ERROR: %s", err.Error())
						} else {
							s.Resp.RespQ <- money
						}
					}
				}(ch)
			}
		}
	}()
}

func (s *TransactionService) WithdrawalMoney(ctx context.Context, request model.TransactionReq) (*model.Response, error) {
	//get balance from db
	balance, err := s.repo.GetBalance(ctx, request.ClientID)
	if err != nil {
		_ = s.repo.ChangeTransactionStatus(ctx, request.TransactionID, TransactionERROR)

		if err.Error() == "client does not exists" {
			s.logger.Info("make new transaction",
				zap.Any("transaction_id", request.TransactionID))

			return &model.Response{
				Result: "",
				Error:  "client does not exists",
			}, err
		}
		return nil, errors.Wrap(err, "get balance from db")
	}

	if balance < 0 || balance-request.Amount < 0 {
		return nil, errors.New("not enough money")
	}

	//decrease balance from client for his request
	err = s.repo.BalanceDecrease(ctx, request)
	if err != nil {
		_ = s.repo.ChangeTransactionStatus(ctx, request.TransactionID, TransactionERROR)
		return nil, errors.Wrap(err, "decrease balance")
	}

	//end transaction
	err = s.repo.ChangeTransactionStatus(ctx, request.TransactionID, TransactionEnd)
	if err != nil {
		_ = s.repo.ChangeTransactionStatus(ctx, request.TransactionID, TransactionERROR)
		return nil, errors.Wrap(err, "transaction end")
	}

	s.logger.Info("make new transaction",
		zap.Any("transaction_id", request.TransactionID))

	return &model.Response{Result: model.ResultOK, Error: ""}, nil
}

func (s *TransactionService) StartTransaction(ctx context.Context, request model.TransactionReq) (string, error) {
	balance, err := s.repo.GetBalance(ctx, request.ClientID)
	if err != nil {
		return "", errors.Wrap(err, "get balance from db")
	}

	transactionID := GetTransactionID()
	trans := &model.Transaction{
		TransactionID:            transactionID,
		ClientID:                 request.ClientID,
		BalanceBeforeTransaction: balance,
		WithdrawalAmount:         request.Amount,
		TransactionStatus:        TransactionStart,
	}

	err = s.repo.TransactionStart(ctx, trans)
	if err != nil {
		return "", errors.Wrap(err, "start transaction")
	}

	return transactionID, nil
}

func (s *TransactionService) RecoveryUncompletedTransactions() error {
	ctx := context.Background()
	transactions, err := s.repo.GetUncompletedTransactions(ctx, TransactionStart)
	if err != nil {
		return errors.Wrap(err, "get uncompleted transaction")
	}

	for _, val := range transactions {
		balance, err := s.repo.GetBalance(ctx, val.ClientID)
		if err != nil {
			return errors.Wrap(err, "get balance for recovery")
		}

		if balance == val.BalanceBeforeTransaction {
			req := model.TransactionReq{
				ClientID: val.ClientID,
				Amount:   val.WithdrawalAmount,
			}

			_, err := s.WithdrawalMoney(ctx, req)
			if err != nil {
				return errors.Wrap(err, "recovery withdrawal")
			}
		}
	}

	return nil
}

func (s *TransactionService) GetAllTransactions(ctx context.Context, req model.RequestGetAllTransaction) (string, error) {
	transactions, err := s.repo.GetAllTransactions(ctx, req.ClientID)
	if err != nil {
		return "", errors.Wrap(err, "get all transactions")
	}

	var answerClient string
	for _, val := range transactions {
		answerClient += fmt.Sprintf(
			"TRANSACTION ID = %s  AMOUNT REQUEST = %d  STATUS TRANSACTION = %s   ",
			val.TransactionID,
			val.WithdrawalAmount,
			val.TransactionStatus)
	}

	s.logger.Info("make status request",
		zap.Any("client_id", req.ClientID))

	return answerClient, nil
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
