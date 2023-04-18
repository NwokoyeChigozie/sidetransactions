package external_models

type DebitWalletRequest struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	BusinessID    int     `json:"business_id"`
	EscrowWallet  string  `json:"escrow_wallet"`
	TransactionID string  `json:"transaction_id"`
}

type CreditWalletRequest struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	BusinessID    int     `json:"business_id"`
	IsRefund      bool    `json:"is_refund"`
	EscrowWallet  string  `json:"escrow_wallet"`
	TransactionID string  `json:"transaction_id"`
}

type WalletBalance struct {
	ID        uint    `json:"id"`
	AccountID int     `json:"account_id"`
	Available float64 `json:"available"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Currency  string  `json:"currency"`
}

type WalletBalanceResponse struct {
	Status  string        `json:"status"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    WalletBalance `json:"data"`
}

type WalletTransferRequest struct {
	SenderAccountID    int     `json:"sender_account_id"`
	RecipientAccountID int     `json:"recipient_account_id"`
	InitialAmount      float64 `json:"initial_amount"`
	FinalAmount        float64 `json:"final_amount"`
	RateID             int     `json:"rate_id"`
	SenderCurrency     string  `json:"sender_currency"`
	RecipientCurrency  string  `json:"recipient_currency"`
	TransactionID      string  `json:"transaction_id"`
	Refund             bool    `json:"refund"`
}
type WalletTransferResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
