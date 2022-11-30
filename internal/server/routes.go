package server

import (
	"github.com/BlackRRR/transaction-service/internal/model"
	"github.com/BlackRRR/transaction-service/internal/services"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	transaction *services.Reader
	queue       *model.Queue
	response    *model.Writer
}

func NewController(service *services.Reader) *Controller {
	return &Controller{
		transaction: service,
		queue:       service.ReadQ,
		response:    service.Resp,
	}

}

func (c *Controller) InitRoutes() *gin.Engine {
	g := gin.New()

	//pen for get money
	g.GET("/create-transaction/get-money", c.GetMoney)

	return g
}
