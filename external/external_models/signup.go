package external_models

type CreateUserRequestModel struct {
	BusinessID      int    `json:"business_id"`
	EmailAddress    string `json:"email_address"`
	PhoneNumber     string `json:"phone_number"`
	AccountType     string `json:"account_type"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Username        string `json:"username"`
	ReferralCode    string `json:"referral_code"`
	Password        string `json:"password"`
	Country         string `json:"country"`
	WebhookURI      string `json:"webhook_uri"`
	BusinessName    string `json:"business_name"`
	BusinessType    string `json:"business_type"`
	BusinessAddress string `json:"business_address"`
}

type SignupModel struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}
