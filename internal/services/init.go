package services

import (
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/BlackRRR/transaction-service/internal/repository"
	"go.uber.org/zap"
)

type TransactionService struct {
	logger *zap.Logger
	repo   repository.Implementation
	ReadQ  *model.Queue
	Resp   *model.Writer
}

func NewReader(logger *zap.Logger, impl repository.Implementation) *TransactionService {
	return &TransactionService{
		logger: logger,
		repo:   impl,
		ReadQ:  &model.Queue{UserQ: make(map[int64]chan model.TransactionReq)},
		Resp:   &model.Writer{RespQ: make(chan *model.Response)},
	}
}
