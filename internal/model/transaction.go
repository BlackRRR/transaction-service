package model

type Transaction struct {
	TransactionID            string `json:"transaction_id"`
	ClientID                 int64  `json:"client_id"`
	BalanceBeforeTransaction int64  `json:"balance_before"`
	BalanceAfterTransaction  int64  `json:"balance_after"`
	WithdrawalAmount         int64  `json:"withdrawal_amount"`
	TransactionStatus        string `json:"transaction_status"`
}
