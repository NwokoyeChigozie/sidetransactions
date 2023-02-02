package external_models

type Payment struct {
	ID               int64   `json:"id"`
	PaymentID        string  `json:"payment_id"`
	TransactionID    string  `json:"transaction_id"`
	TotalAmount      float64 `json:"total_amount"`
	EscrowCharge     float64 `json:"escrow_charge"`
	IsPaid           bool    `json:"is_paid"`
	PaymentMadeAt    string  `json:"payment_made_at"`
	DeletedAt        string  `json:"deleted_at"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	AccountID        int64   `json:"account_id"`
	BusinessID       int64   `json:"business_id"`
	Currency         string  `json:"currency"`
	ShippingFee      float64 `json:"shipping_fee"`
	DisburseCurrency string  `json:"disburse_currency"`
	PaymentType      string  `json:"payment_type"`
	BrokerCharge     float64 `json:"broker_charge"`
}

type CreatePaymentRequestWithToken struct {
	TransactionID string  `json:"transaction_id" `
	TotalAmount   float64 `json:"total_amount"`
	ShippingFee   float64 `json:"shipping_fee"`
	BrokerCharge  float64 `json:"broker_charge"`
	EscrowCharge  float64 `json:"escrow_charge"`
	Currency      string  `json:"currency"`
	Token         string  `json:"token"`
}

type CreatePaymentRequest struct {
	TransactionID string  `json:"transaction_id" `
	TotalAmount   float64 `json:"total_amount"`
	ShippingFee   float64 `json:"shipping_fee"`
	BrokerCharge  float64 `json:"broker_charge"`
	EscrowCharge  float64 `json:"escrow_charge"`
	Currency      string  `json:"currency"`
}

type CreatePaymentResponse struct {
	Status  string  `json:"status"`
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    Payment `json:"data"`
}
