package external_models

type EmailNotificationRequest struct {
	EmailAddress string `json:"email_address"`
	AccountId    uint   `json:"account_id"`
	Code         uint   `json:"code"`
	Token        string `json:"token"`
}
type AccountIDRequestModel struct {
	AccountId uint `json:"account_id"`
}
type TransactionIDRequestModel struct {
	TransactionId string `json:"transaction_id"`
}
type TransactionIDAccountIDRequestModel struct {
	TransactionId string `json:"transaction_id"`
	AccountId     uint   `json:"account_id"`
}
type DueDateExtensionProposalRequestModel struct {
	TransactionId string `json:"transaction_id"`
	Note          string `json:"note"`
}
