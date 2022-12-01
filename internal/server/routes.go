package server

import (
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/BlackRRR/transaction-service/internal/services"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	transaction *services.TransactionService
	queue       *model.Queue
	response    *model.Writer
}

func NewController(service *services.TransactionService) *Controller {
	return &Controller{
		transaction: service,
		queue:       service.ReadQ,
		response:    service.Resp,
	}

}

func (c *Controller) InitRoutes() *gin.Engine {
	g := gin.New()

	//route for get money
	g.GET("/create-transaction/get-money", c.GetMoney)
	g.GET("/transaction-status", c.StatusTransaction)

	return g
}
