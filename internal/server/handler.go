package server

import (
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

func (c *Controller) GetMoney(gc *gin.Context) {
	var req model.TransactionReq
	//get request
	if err := gc.BindJSON(&req); err != nil {
		NewLinksError(gc, http.StatusBadRequest, "invalid req body")
		return
	}

	transactionID, err := c.transaction.StartTransaction(gc, req)
	if err != nil {
		if strings.Contains(err.Error(), "client does not exist") {
			NewLinksError(gc, http.StatusInternalServerError, "CLIENT DOES NOT EXISTS")
			return
		}
		NewLinksError(gc, http.StatusInternalServerError, "CANCEL TRANSACTION, INTERNAL SERVER ERROR")
		return
	}

	tr := model.TransactionReq{
		TransactionID: transactionID,
		ClientID:      req.ClientID,
		Amount:        req.Amount,
	}

	//check if request exitst for this user
	_, exist := c.queue.UserQ[req.ClientID]
	if !exist {
		c.queue.UserQ[req.ClientID] = make(chan model.TransactionReq)
	}

	//send request to channel
	c.queue.UserQ[req.ClientID] <- tr

	var resp *model.Response

	wg := &sync.WaitGroup{}
	wg.Add(1)

	//waiting for response from service
	go func(wg *sync.WaitGroup) {
		for {
			select {
			case resp = <-c.response.RespQ:
				wg.Done()
				return
			}
		}
	}(wg)

	wg.Wait()
	if resp.Error != "" {
		NewLinksError(gc, http.StatusInternalServerError, resp.Error)
		return
	}

	gc.JSON(http.StatusOK, resp.Result)
}

func (c *Controller) StatusTransaction(gc *gin.Context) {
	var req model.RequestGetAllTransaction

	//get request
	if err := gc.BindJSON(&req); err != nil {
		NewLinksError(gc, http.StatusBadRequest, "invalid req body")
		return
	}

	//get all transactions info
	answer, err := c.transaction.GetAllTransactions(gc, req)
	if err != nil {
		NewLinksError(gc, http.StatusInternalServerError, "CANCEL REQUEST, INTERNAL SERVER ERROR")
		return
	}

	if answer == "" {
		gc.JSON(http.StatusOK, "NO AVAILABLE TRANSACTIONS")
	}

	gc.JSON(http.StatusOK, answer)
}

func NewLinksError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, model.ErrorResponse{ErrorMessage: message})
}
