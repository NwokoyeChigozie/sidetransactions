package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func ListTransactionsByIDService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string) (models.TransactionByIDResponse, int, error) {
	var (
		transaction        = models.Transaction{TransactionID: transactionID}
		productTransaction = models.ProductTransaction{TransactionID: transactionID}
		parties            = []models.TransactionParty{}
		transactionFile    = models.TransactionFile{TransactionID: transactionID}
	)

	transactions, err := transaction.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, code, fmt.Errorf("transaction not found: %v", err.Error())
	}

	productTransactions, err := productTransaction.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	if transaction.Type == "milestone" {
		transactionParty := models.TransactionParty{TransactionPartiesID: transactions[0].PartiesID}
		parties, err = transactionParty.GetAllByTransactionPartiesID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
	} else {
		transactionParty := models.TransactionParty{TransactionID: transaction.TransactionID}
		parties, err = transactionParty.GetAllByTransactionID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
	}

	transactionFiles, err := transactionFile.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	transactionReponse := resolveTransactionAndListTransactionResponse(transaction)
	transactionReponse.Products = productTransactions

	partiess, members, err := getPartiesAndMembersFromParties(extReq, parties)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	if len(parties) > 0 {
		transactionReponse.Parties = partiess
		transactionReponse.Members = members
	} else {
		transactionParty := models.TransactionParty{TransactionPartiesID: transactions[0].PartiesID}
		tParties, err := transactionParty.GetAllByTransactionPartiesID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
		partiess, members, err = getPartiesAndMembersFromParties(extReq, tParties)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}

		transactionReponse.Parties = partiess
		transactionReponse.Members = members
	}

	transactionReponse.Files = transactionFiles
	milestones := map[int][]models.MilestonesResponse{}

	for i, t := range transactions {
		totalAmount, mileSresponse := resolveTransactionForAmountAndMilestoneResponse(extReq, i, t)
		transactionReponse.TotalAmount = totalAmount
		currentArray := milestones[i]
		currentArray = append(currentArray, mileSresponse)
		milestones[i] = currentArray

		otherTransactions, err := t.GetAllOthersByIDAndPartiesID(db.Transaction)
		if err != nil {
			logger.Info("error getting other transactions", err.Error())
		}

		for oi, ot := range otherTransactions {
			_, otherMileSresponse := resolveTransactionForAmountAndMilestoneResponse(extReq, oi, ot)
			currentArray := milestones[i]
			currentArray = append(currentArray, otherMileSresponse)
			milestones[i] = currentArray
		}

	}

	transactionBroker := models.TransactionBroker{TransactionID: transactionID}
	activity := models.ActivityLog{TransactionID: transactionID}

	code, err = transactionBroker.GetTransactionBrokerByTransactionID(db.Transaction)
	if err != nil && code == http.StatusInternalServerError {
		logger.Info("error getting transaction broker", err.Error())
	}

	activities, err := activity.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	country, err := GetCountryByNameOrCode(extReq, logger, transaction.Country)
	if err != nil && code == http.StatusInternalServerError {
		logger.Info("error getting country", err.Error())
	}

	dDateFormatted, err := utility.FormatDate(transaction.DueDate, "2006-01-02", "2006-01-02 15:04:05")
	if err != nil {
		dDateFormatted = ""
	}

	transactionState := models.TransactionState{TransactionID: transactionID, Status: "Closed"}
	code, err = transactionState.GetTransactionStateByTransactionIDAndStatus(db.Transaction)
	if err != nil && code == http.StatusInternalServerError {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	transactionDispute := models.TransactionDispute{TransactionID: transactionID}
	isDisputed, err := transactionDispute.IsDisputed(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	transactionReponse.Milestones = milestones[0]
	transactionReponse.Broker = transactionBroker
	transactionReponse.Activities = activities
	transactionReponse.Country = country
	transactionReponse.DueDateFormatted = dDateFormatted
	transactionReponse.TransactionClosedAt = transactionState.CreatedAt
	transactionReponse.IsDisputed = isDisputed

	return transactionReponse, http.StatusOK, nil
}

func ListTransactionsByUssdCodeService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, ussdCode int) (models.TransactionByIDResponse, int, error) {
	var (
		transaction = models.Transaction{TransUssdCode: ussdCode}
	)

	code, err := transaction.GetTransactionByUssdCode(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, code, err
	}

	return ListTransactionsByIDService(extReq, logger, db, transaction.TransactionID)
}

func ListTransactionsService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListTransactionsRequest, paginator postgresql.Pagination) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transaction  = models.Transaction{}
		transactions = []models.Transaction{}
	)

	if req.Status != "" {
		transaction.Status = req.Status
	} else {
		transaction.Status = GetTransactionStatus(req.StatusCode)
	}

	transactions, pagination, err := transaction.GetAllByAndQueries(db.Transaction, false, req.Filter, "id", "asc", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	var transactionsResponses []models.TransactionByIDResponse

	for _, t := range transactions {
		transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, t.TransactionID)
		if err != nil {
			logger.Error("list transaction by id error", err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, transactionResponse)
		}

	}

	return transactionsResponses, pagination, http.StatusOK, nil

}

func ListTransactionsByBusinessService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListTransactionByBusinessRequest, paginator postgresql.Pagination) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transaction  = models.Transaction{BusinessID: req.BusinessID, IsPaylinked: req.Paylinked}
		transactions = []models.Transaction{}
	)

	if req.Status != "" {
		transaction.Status = req.Status
	} else {
		transaction.Status = GetTransactionStatus(req.StatusCode)
	}

	transactions, pagination, err := transaction.GetAllByAndQueries(db.Transaction, req.Paylinked, req.Filter, "id", "asc", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	var transactionsResponses []models.TransactionByIDResponse

	for _, t := range transactions {
		transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, t.TransactionID)
		if err != nil {
			logger.Error("list transaction by id error", err.Error())
		} else {
			payment, err := ListPayment(extReq, t.TransactionID)
			if err != nil {
				transactionResponse.TotalAmount = 0
				transactionResponse.EscrowCharge = 0
				logger.Error("list payment by transaction id error", err.Error())
			} else {
				transactionResponse.TotalAmount = payment.TotalAmount
				transactionResponse.EscrowCharge = payment.EscrowCharge
			}

			transactionsResponses = append(transactionsResponses, transactionResponse)
		}

	}

	return transactionsResponses, pagination, http.StatusOK, nil

}
func ListByBusinessFromMondayToThursdayService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListByBusinessFromMondayToThursdayRequest, paginator postgresql.Pagination) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transaction  = models.Transaction{BusinessID: req.BusinessID, IsPaylinked: req.Paylinked}
		transactions = []models.Transaction{}
	)

	if req.Status != "" {
		transaction.Status = req.Status
	} else {
		transaction.Status = GetTransactionStatus(req.StatusCode)
	}

	transactions, pagination, err := transaction.GetAllByAndQueries(db.Transaction, req.Paylinked, "monday_to_thursday", "id", "asc", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	var transactionsResponses []models.TransactionByIDResponse

	for _, t := range transactions {
		transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, t.TransactionID)
		if err != nil {
			logger.Error("list transaction by id error", err.Error())
		} else {
			payment, err := ListPayment(extReq, t.TransactionID)
			if err != nil {
				transactionResponse.TotalAmount = 0
				transactionResponse.EscrowCharge = 0
				logger.Error("list payment by transaction id error", err.Error())
			} else {
				transactionResponse.TotalAmount = payment.TotalAmount
				transactionResponse.EscrowCharge = payment.EscrowCharge
			}

			transactionsResponses = append(transactionsResponses, transactionResponse)
		}

	}

	return transactionsResponses, pagination, http.StatusOK, nil

}

func ListTransactionsByUserService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListTransactionByUserRequest, paginator postgresql.Pagination, user external_models.User) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transactions          = []models.Transaction{}
		transactionsResponses = []models.TransactionByIDResponse{}
		transactionParty      = models.TransactionParty{AccountID: int(user.AccountID), Role: req.Role}
	)

	transactionParties, pagination, err := transactionParty.GetAllByAndQueriesForUniqueValue(db.Transaction, "", "id", "desc", "transaction_id", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	for _, tp := range transactionParties {
		lTransaction := models.Transaction{TransactionID: tp.TransactionID, IsPaylinked: req.Paylinked, Status: GetTransactionStatus(req.StatusCode)}
		code, err := lTransaction.GetLatestByAndQueries(db.Transaction, req.Paylinked, "")
		if err != nil {
			if code == http.StatusInternalServerError {
				return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, code, err
			}
			logger.Error("error getting lastest transaction byy transaction id", err.Error())
		} else {
			transactions = append(transactions, lTransaction)
		}
	}

	for _, t := range transactions {
		transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, t.TransactionID)
		if err != nil {
			logger.Error("list transaction by id error", err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, transactionResponse)
		}

	}

	return transactionsResponses, pagination, http.StatusOK, nil

}

func ListArchivedTransactionsService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, paginator postgresql.Pagination, user external_models.User) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transactions          = []models.Transaction{}
		transactionsResponses = []models.TransactionByIDResponse{}
		transactionParty      = models.TransactionParty{AccountID: int(user.AccountID), Role: "sender"}
	)

	transactionParties, pagination, err := transactionParty.GetAllByAndQueriesForUniqueValue(db.Transaction, "", "id", "desc", "transaction_id", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	for _, tp := range transactionParties {
		lTransaction := models.Transaction{TransactionID: tp.TransactionID, Status: "Deleted"}
		code, err := lTransaction.GetLatestByAndQueries(db.Transaction, false, "")
		if err != nil {
			if code == http.StatusInternalServerError {
				return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, code, err
			}
			logger.Error("error getting lastest transaction byy transaction id", err.Error())
		} else {
			transactions = append(transactions, lTransaction)
		}
	}

	for _, t := range transactions {
		transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, t.TransactionID)
		if err != nil {
			logger.Error("list transaction by id error", err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, transactionResponse)
		}

	}

	return transactionsResponses, pagination, http.StatusOK, nil

}

func resolveTransactionForAmountAndMilestoneResponse(extReq request.ExternalRequest, i int, t models.Transaction) (float64, models.MilestonesResponse) {
	var (
		totalAmount float64 = 0
	)
	var (
		titleSlice = strings.Split(t.Title, ";")
		title      = ""
		index      = 1
	)

	if len(titleSlice) >= 3 {
		totalAmount, _ = strconv.ParseFloat(titleSlice[2], 64)
	}

	if len(titleSlice) > 0 {
		title = titleSlice[1]
	}

	if len(titleSlice) >= 4 {
		index, _ = strconv.Atoi(titleSlice[3])
		if index == 0 {
			index = 1
		}
	}

	var recipients []models.MileStoneRecipient
	var recipientsResponse []models.MilestonesRecipientResponse

	err := json.Unmarshal([]byte(t.Recipients), &recipients)
	if err != nil {
		extReq.Logger.Info("error unmarshaling recipient json string to struct", t.Recipients, err.Error())
	}

	for _, r := range recipients {
		user, _ := GetUserWithAccountID(extReq, r.AccountID)
		accountName := ""
		if user.ID != 0 {
			accountName = user.Lastname + " " + user.Firstname
		}

		recipientsResponse = append(recipientsResponse, models.MilestonesRecipientResponse{
			AccountID:   r.AccountID,
			AccountName: accountName,
			Email:       user.EmailAddress,
			PhoneNumber: user.PhoneNumber,
			Amount:      r.Amount,
		})

	}

	return totalAmount, models.MilestonesResponse{
		Index:            index,
		MilestoneID:      t.MilestoneID,
		Title:            title,
		Amount:           t.Amount,
		Status:           t.Status,
		InspectionPeriod: t.InspectionPeriod,
		DueDate:          t.DueDate,
		Recipients:       recipientsResponse,
	}
}

func getPartiesAndMembersFromParties(extReq request.ExternalRequest, parties []models.TransactionParty) (map[string]models.TransactionParty, []models.PartyResponse, error) {
	var (
		partiess = map[string]models.TransactionParty{}
		members  = []models.PartyResponse{}
	)
	for _, p := range parties {
		user, _ := GetUserWithAccountID(extReq, p.AccountID)
		partiess[p.Role] = p
		accountName := ""
		if user.ID != 0 {
			accountName = user.Lastname + " " + user.Firstname
		}
		var roleCapabilities models.PartyAccessLevel

		inrec, err := json.Marshal(p.RoleCapabilities)
		if err != nil {
			return partiess, members, err
		}

		err = json.Unmarshal(inrec, &roleCapabilities)
		if err != nil {
			return partiess, members, err
		}
		members = append(members, models.PartyResponse{
			PartyID:     int(p.ID),
			AccountID:   p.AccountID,
			AccountName: accountName,
			PhoneNumber: user.PhoneNumber,
			Email:       user.EmailAddress,
			Role:        p.Role,
			Status:      p.Status,
			AccessLevel: roleCapabilities,
		})
	}
	return partiess, members, nil
}

func resolveTransactionAndListTransactionResponse(transaction models.Transaction) models.TransactionByIDResponse {
	var rRrecipients []models.MileStoneRecipient
	json.Unmarshal([]byte(transaction.Recipients), &rRrecipients)
	return models.TransactionByIDResponse{
		ID:               transaction.ID,
		TransactionID:    transaction.TransactionID,
		PartiesID:        transaction.PartiesID,
		MilestoneID:      transaction.MilestoneID,
		BrokerID:         transaction.BrokerID,
		Title:            transaction.Title,
		Type:             transaction.Type,
		Description:      transaction.Description,
		Amount:           transaction.Amount,
		Status:           transaction.Status,
		Quantity:         transaction.Quantity,
		InspectionPeriod: transaction.InspectionPeriod,
		DueDate:          transaction.DueDate,
		ShippingFee:      transaction.ShippingFee,
		GracePeriod:      transaction.GracePeriod,
		Currency:         transaction.Currency,
		DeletedAt:        transaction.DeletedAt,
		CreatedAt:        transaction.CreatedAt,
		UpdatedAt:        transaction.UpdatedAt,
		BusinessID:       transaction.BusinessID,
		IsPaylinked:      transaction.IsPaylinked,
		Source:           transaction.Source,
		TransUssdCode:    transaction.TransUssdCode,
		Recipients:       rRrecipients,
		DisputeHandler:   transaction.DisputeHandler,
		AmountPaid:       transaction.AmountPaid,
		EscrowCharge:     transaction.EscrowCharge,
		EscrowWallet:     transaction.EscrowWallet,
	}
}
func resolveTransactionAndListTransactionResponse2(transaction models.Transaction) models.TransactionByIDResponse {
	v := models.TransactionByIDResponse{}
	utility.CopyStruct(&transaction, &v)
	return v
}
