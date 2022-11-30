package model

import (
	"github.com/gin-gonic/gin"
)

const (
	ResultOK  = "OK"
	ResultErr = "SOME ERROR ON SERVER, TRY AGAIN LATER"
)

// Request from client for transaction
type Request struct {
	ClientID      int64  `json:"client_id"`
	Amount        int64  `json:"amount"`
	TransactionID string `json:"transaction_id"`
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// Queue channel for requests from client
type Queue struct {
	UserQ map[int64]chan Request
}

// Writer response channel for answer to client
type Writer struct {
	RespQ chan *Response
}

// custom error response
type errorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func NewLinksError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
