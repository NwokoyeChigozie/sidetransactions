package mocks

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/mocks/appruve_mocks"
	"github.com/vesicash/transactions-ms/external/mocks/auth_mocks"
	"github.com/vesicash/transactions-ms/external/mocks/ipstack_mocks"
	"github.com/vesicash/transactions-ms/external/mocks/monnify_mocks"
	"github.com/vesicash/transactions-ms/external/mocks/notification_mocks"
	"github.com/vesicash/transactions-ms/external/mocks/payment_mocks"
	"github.com/vesicash/transactions-ms/external/mocks/rave_mocks"
	"github.com/vesicash/transactions-ms/utility"
)

type ExternalRequest struct {
	Logger     *utility.Logger
	Test       bool
	RequestObj RequestObj
}

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	RequestData  interface{}
	DecodeMethod string
	Logger       *utility.Logger
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	switch name {
	case "get_user":
		return auth_mocks.GetUser(er.Logger, data)
	case "get_user_credential":
		return auth_mocks.GetUserCredential(er.Logger, data)
	case "create_user_credential":
		return auth_mocks.CreateUserCredential(er.Logger, data)
	case "update_user_credential":
		return auth_mocks.UpdateUserCredential(er.Logger, data)
	case "get_user_profile":
		return auth_mocks.GetUserProfile(er.Logger, data)
	case "get_business_profile":
		return auth_mocks.GetBusinessProfile(er.Logger, data)
	case "get_country":
		return auth_mocks.GetCountry(er.Logger, data)
	case "get_bank_details":
		return auth_mocks.GetBankDetails(er.Logger, data)
	case "get_access_token":
		return auth_mocks.GetAccessToken(er.Logger)
	case "validate_on_auth":
		return auth_mocks.ValidateOnAuth(er.Logger, data)
	case "validate_authorization":
		return auth_mocks.ValidateAuthorization(er.Logger, data)
	case "send_verification_email":
		return notification_mocks.SendVerificationEmail(er.Logger, data)
	case "send_welcome_email":
		return notification_mocks.SendWelcomeEmail(er.Logger, data)
	case "send_email_verified_notification":
		return notification_mocks.SendEmailVerifiedNotification(er.Logger, data)
	case "send_sms_to_phone":
		return notification_mocks.SendSendSMSToPhone(er.Logger, data)
	case "monnify_login":
		return monnify_mocks.MonnifyLogin(er.Logger, data)
	case "monnify_match_bvn_details":
		return monnify_mocks.MonnifyMatchBvnDetails(er.Logger, data)
	case "appruve_verify_id":
		return appruve_mocks.AppruveVerifyID(er.Logger, data)
	case "verification_failed_notification":
		return notification_mocks.VerificationFailedNotification(er.Logger, data)
	case "verification_successful_notification":
		return notification_mocks.VerificationSuccessfulNotification(er.Logger, data)
	case "rave_resolve_bank_account":
		return rave_mocks.RaveResolveBankAccount(er.Logger, data)
	case "ipstack_resolve_ip":
		return ipstack_mocks.IpstackResolveIp(er.Logger, data)
	case "get_authorize":
		return auth_mocks.GetAuthorize(er.Logger, data)
	case "create_authorize":
		return auth_mocks.CreateAuthorize(er.Logger, data)
	case "update_authorize":
		return auth_mocks.UpdateAuthorize(er.Logger, data)
	case "send_authorized_notification":
		return notification_mocks.SendAuthorizedNotification(er.Logger, data)
	case "send_authorization_notification":
		return notification_mocks.SendAuthorizationNotification(er.Logger, data)
	case "set_user_authorization_required_status":
		return auth_mocks.SetUserAuthorizationRequiredStatus(er.Logger, data)
	case "get_business_charge":
		return auth_mocks.GetBusinessCharge(er.Logger, data)
	case "init_business_charge":
		return auth_mocks.InitBusinessCharge(er.Logger, data)
	case "create_payment":
		return payment_mocks.CreatePayment(er.Logger, data)
	case "list_payment":
		return payment_mocks.ListPayment(er.Logger, data)
	case "signup_user":
		return auth_mocks.SignupUser(er.Logger, data)
	case "send_new_transaction_notification":
		return notification_mocks.SendNewTransactionNotification(er.Logger, data)
	case "send_transaction_accepted_notification":
		return notification_mocks.SendTransactionAcceptedNotification(er.Logger, data)
	case "send_transaction_rejected_notification":
		return notification_mocks.SendTransactionRejectedNotification(er.Logger, data)
	case "send_transaction_delivered_rejected_notification":
		return notification_mocks.SendTransactionDeliveredRejectedNotification(er.Logger, data)
	case "send_dispute_opened_notification":
		return notification_mocks.SendDisputeOpenedNotification(er.Logger, data)
	case "send_transaction_delivered_notification":
		return notification_mocks.SendTransactionDeliveredNotification(er.Logger, data)
	case "send_due_date_proposal_notification":
		return notification_mocks.SendDueDateProposalNotification(er.Logger, data)
	case "send_due_date_extended_notification":
		return notification_mocks.SendDueDateExtendedNotification(er.Logger, data)
	case "send_transaction_delivered_accepted_notification":
		return notification_mocks.SendTransactionDeliveredAcceptedNotification(er.Logger, data)
	case "get_access_token_by_key":
		return auth_mocks.GetAccessTokenByKey(er.Logger, data)
	case "request_manual_refund":
		return payment_mocks.RequestManualRefund(er.Logger, data)
	case "wallet_transfer":
		return payment_mocks.WalletTransfer(er.Logger, data)
	case "debit_wallet":
		return payment_mocks.DebitWallet(er.Logger, data)
	case "credit_wallet":
		return payment_mocks.CreditWallet(er.Logger, data)
	default:
		return nil, fmt.Errorf("request not found")
	}
}
