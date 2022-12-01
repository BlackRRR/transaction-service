package model

const (
	ResultOK  = "OK"
	ResultErr = "SOME ERROR ON SERVER, TRY AGAIN LATER"
)

// Request from client for transaction
type Request struct {
	ClientID int64 `json:"client_id"`
	Amount   int64 `json:"amount"`
}

// TransactionReq created for transaction from request
type TransactionReq struct {
	TransactionID string `json:"transaction_id"`
	ClientID      int64  `json:"client_id"`
	Amount        int64  `json:"amount"`
}

// Response to client
type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// RequestGetAllTransaction from server
type RequestGetAllTransaction struct {
	ClientID int64 `json:"client_id"`
}
