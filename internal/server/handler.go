package server

import (
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func (c *Controller) GetMoney(gc *gin.Context) {
	var req model.Request

	//get request
	if err := gc.BindJSON(&req); err != nil {
		model.NewLinksError(gc, http.StatusBadRequest, "invalid req body")
		return
	}

	//start transaction
	transactionID, err := c.transaction.StartTransaction(gc, req)
	if err != nil {
		model.NewLinksError(gc, http.StatusInternalServerError, err.Error())
		return
	}

	req.TransactionID = transactionID

	//check if request exitst for this user
	_, exist := c.queue.UserQ[req.ClientID]
	if !exist {
		c.queue.UserQ[req.ClientID] = make(chan model.Request)
	}

	//send request to channel
	c.queue.UserQ[req.ClientID] <- req

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
		model.NewLinksError(gc, http.StatusInternalServerError, resp.Error)
		return
	}

	gc.JSON(http.StatusOK, resp.Result)
}
