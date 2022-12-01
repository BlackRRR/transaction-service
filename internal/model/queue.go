package model

// Queue map channel for requests from client
type Queue struct {
	UserQ map[int64]chan TransactionReq
}

// Writer response channel for answer to client
type Writer struct {
	RespQ chan *Response
}
