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

func CreateTransactionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.CreateTransactionRequest, user external_models.User) (models.TransactionCreateResponse, int, error) {
	var (
		transaction               = models.Transaction{}
		businessID                = req.BusinessID
		businessCharge            = external_models.BusinessCharge{}
		transactionID             = utility.RandomString(20)
		transactionPartiesID      = utility.RandomString(20)
		transactionTitle          = req.Title
		transactionDescription    = req.Description
		transactioQuantity        = req.Quantity
		transactionAmount         = req.Amount
		inspectionPeriod          = req.InspectionPeriod
		transactionShippingFee    = req.ShippingFee
		transactionType           = req.Type
		transactionCurrency       = strings.ToUpper(req.Currency)
		transactionPaylinked      = req.Paylinked
		transactionSource         = req.Source
		transactionDisputeHandler = req.DisputeHandler
	)
	if businessID == 0 {
		businessID = int(user.AccountID)
	}

	if transactionDescription == "" {
		transactionDescription = req.Title
	}

	if transactionSource == "" {
		transactionSource = "api"
	}

	dueDate, err := validateDueDate(req.DueDate)
	if err != nil {
		return models.TransactionCreateResponse{}, http.StatusBadRequest, fmt.Errorf("incorrect due date format, try 2006-01-15")
	}
	transactionDueDate, _ := utility.GetUnixString(dueDate, "2006-01-02", "2006-01-02")
	if err := validatePartiesAndMilestones(transactionType, req.Parties, req.Milestones); err != nil {
		return models.TransactionCreateResponse{}, http.StatusBadRequest, err
	}

	gracePeriod, err := validateDueDate(req.GracePeriod)
	if err != nil {
		return models.TransactionCreateResponse{}, http.StatusBadRequest, fmt.Errorf("incorrect grace period format, try 2006-01-15")
	}

	transactionGracePeriod, _ := utility.GetUnixString(gracePeriod, "2006-01-02", "2006-01-02")

	businessCharge, err = getBusinessChargeWithBusinessIDAndCurrency(extReq, businessID, req.Currency)
	if err != nil {
		businessCharge, err = initBusinessCharge(extReq, businessID, req.Currency)
		if err != nil {
			return models.TransactionCreateResponse{}, http.StatusInternalServerError, err
		}
	}

	transactionCountry := businessCharge.Country
	transactionStatus := GetTransactionStatus("draft")

	transactionFiles, err := resolveTransactionFiles(req.Files, transactionID, businessID, db)
	if err != nil {
		return models.TransactionCreateResponse{}, http.StatusInternalServerError, err
	}

	_, partiesResponse, err := resolveParties(extReq, req.Parties, transactionPartiesID, transactionID, db)
	if err != nil {
		return models.TransactionCreateResponse{}, http.StatusInternalServerError, err
	}

	mileStoneResponse := []models.MilestonesResponse{}
	transactionObj := models.ResolveTransactionObj{
		TransactionID:        transactionID,
		TransactionPartiesID: transactionPartiesID,
		Title:                transactionTitle,
		Type:                 transactionType,
		Description:          transactionDescription,
		Amount:               transactionAmount,
		Quantity:             transactioQuantity,
		ShippingFee:          transactionShippingFee,
		GracePeriod:          transactionGracePeriod,
		Currency:             transactionCurrency,
		Country:              transaction.Country,
		BusinessID:           businessID,
		DisputeHandler:       transactionDisputeHandler,
		EscrowWallet:         req.EscrowWallet,
	}
	if transactionType == "oneoff" {
		transaction, mileStoneResponse, err = resolveCreateOneOffTransaction(extReq, req.Milestones, transactionAmount, transactionObj, db)
		if err != nil {
			return models.TransactionCreateResponse{}, http.StatusInternalServerError, err
		}
	} else if transactionType == "milestone" {
		transaction, mileStoneResponse, err = resolveCreateMilestoneTransaction(extReq, req.Milestones, transactionAmount, transactionObj, db)
		if err != nil {
			return models.TransactionCreateResponse{}, http.StatusInternalServerError, err
		}
	} else {
		return models.TransactionCreateResponse{}, http.StatusBadRequest, fmt.Errorf("transaction type not implemented")
	}

	transaction.IsPaylinked = transactionPaylinked
	transaction.Source = transactionSource
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return models.TransactionCreateResponse{}, http.StatusInternalServerError, err
	}
	escrowCharge := getEscrowCharge(businessCharge, getTotalAmoutForMilestones(req.Milestones))
	if transactionSource == "transfer" {
		escrowCharge = 2
	}

	createPaymentPayload := external_models.CreatePaymentRequestWithToken{
		TransactionID: transactionID,
		TotalAmount:   transactionAmount,
		ShippingFee:   transactionShippingFee,
		BrokerCharge:  0,
		EscrowCharge:  escrowCharge,
		Currency:      transactionCurrency,
		Token:         models.Token,
	}

	_, err = CreatePayment(extReq, createPaymentPayload)
	if err != nil {
		return models.TransactionCreateResponse{}, http.StatusInternalServerError, fmt.Errorf("payment creation failed: %v", err)
	}

	activityLog := models.ActivityLog{
		TransactionID: transactionID,
		Description:   "Transaction details have been sent to all invited parties",
	}
	err = activityLog.CreateActivityLog(db.Transaction)
	if err != nil {
		return models.TransactionCreateResponse{}, http.StatusInternalServerError, err
	}

	var rRrecipients []models.MileStoneRecipient
	json.Unmarshal([]byte(transaction.Recipients), &rRrecipients)

	return models.TransactionCreateResponse{
		ID:               transaction.ID,
		TransactionID:    transaction.TransactionID,
		PartiesID:        transaction.PartiesID,
		MilestoneID:      transaction.MilestoneID,
		BrokerID:         transaction.BrokerID,
		Title:            transactionTitle,
		Type:             transaction.Type,
		Description:      transactionDescription,
		Amount:           transactionAmount,
		Status:           transactionStatus,
		Quantity:         transaction.Quantity,
		InspectionPeriod: strconv.Itoa(inspectionPeriod),
		DueDate:          transactionDueDate,
		ShippingFee:      transactionShippingFee,
		Currency:         transactionCurrency,
		IsPaylinked:      transaction.IsPaylinked,
		Source:           transactionSource,
		TransUssdCode:    transaction.TransUssdCode,
		Recipients:       rRrecipients,
		DisputeHandler:   transactionDisputeHandler,
		EscrowCharge:     transaction.EscrowCharge,
		EscrowWallet:     transaction.EscrowWallet,
		Country:          transactionCountry,
		Parties:          partiesResponse,
		Files:            transactionFiles,
		Milestones:       mileStoneResponse,
	}, http.StatusOK, nil
}

func getEscrowCharge(businessCharge external_models.BusinessCharge, totalAmountForMilestones float64) float64 {
	var (
		charge float64 = 0
	)

	bCharge, err := strconv.ParseFloat(businessCharge.BusinessCharge, 64)
	if err != nil {
		bCharge = 0
	}
	vCharge, err := strconv.ParseFloat(businessCharge.VesicashCharge, 64)
	if err != nil {
		vCharge = 0
	}
	processingFee, err := strconv.ParseFloat(businessCharge.ProcessingFee, 64)
	if err != nil {
		processingFee = 0
	}

	if businessCharge.ChargeMin != nil && businessCharge.ChargeMid != nil && businessCharge.ChargeMax != nil {
		chargeMin := utility.ConvertStringInterfaceToStringFloat(*businessCharge.ChargeMin)
		chargeMid := utility.ConvertStringInterfaceToStringFloat(*businessCharge.ChargeMid)
		chargeMax := utility.ConvertStringInterfaceToStringFloat(*businessCharge.ChargeMax)
		chargeMinAmount, ok1 := chargeMin["amount"]
		chargeMinCharge, ok11 := chargeMin["charge"]
		chargeMidAmount, ok2 := chargeMid["amount"]
		chargeMidCharge, ok22 := chargeMid["charge"]
		_, ok3 := chargeMax["amount"]
		chargeMaxCharge, ok33 := chargeMax["charge"]

		if ok1 && ok11 && ok2 && ok22 && ok3 && ok33 {
			if totalAmountForMilestones <= chargeMinAmount {
				charge = chargeMinCharge
			} else if totalAmountForMilestones >= chargeMinAmount && totalAmountForMilestones <= chargeMidAmount {
				charge = chargeMidCharge
			} else if totalAmountForMilestones >= chargeMidAmount {
				charge = chargeMaxCharge
			}

		} else {
			charge = utility.PercentageOf(totalAmountForMilestones, bCharge+vCharge) + processingFee
		}

	} else {
		charge = utility.PercentageOf(totalAmountForMilestones, bCharge+vCharge) + processingFee
	}
	return charge
}

func resolveCreateOneOffTransaction(extReq request.ExternalRequest, milestones []models.MileStone, transactionAmount float64, transactionObj models.ResolveTransactionObj, db postgresql.Databases) (models.Transaction, []models.MilestonesResponse, error) {
	var (
		escrowCharge       = transactionAmount - getTotalAmoutForMilestones(milestones)
		milestonesResponse = []models.MilestonesResponse{}
		transactionM       = models.Transaction{}
	)

	for i, m := range milestones {
		milestoneID := utility.RandomString(20)
		transUssdCode := utility.GetRandomNumbersInRange(10000, 99999)
		dueDate, _ := utility.GetUnixString(m.DueDate, "2006-01-02", "2006-01-02")

		recipientsJson, err := json.Marshal(m.Recipients)
		if err != nil {
			return transactionM, milestonesResponse, err
		}
		transaction := models.Transaction{
			TransactionID:    transactionObj.TransactionID,
			PartiesID:        transactionObj.TransactionPartiesID,
			Title:            transactionObj.Title + ";" + m.Title + ";" + strconv.Itoa(int(transactionObj.Amount)) + ";" + strconv.Itoa(i+1),
			Type:             transactionObj.Type,
			Description:      transactionObj.Description,
			MilestoneID:      milestoneID,
			Amount:           transactionObj.Amount,
			Status:           m.Status,
			Quantity:         transactionObj.Quantity,
			InspectionPeriod: strconv.Itoa(m.InspectionPeriod),
			DueDate:          dueDate,
			ShippingFee:      transactionObj.ShippingFee,
			GracePeriod:      transactionObj.GracePeriod,
			Currency:         transactionObj.Currency,
			Country:          transactionObj.Country,
			BusinessID:       transactionObj.BusinessID,
			EscrowCharge:     escrowCharge,
			TransUssdCode:    transUssdCode,
			Recipients:       string(recipientsJson),
			DisputeHandler:   transactionObj.DisputeHandler,
			EscrowWallet:     transactionObj.EscrowWallet,
		}

		err = transaction.CreateTransaction(db.Transaction)
		if err != nil {
			return transactionM, milestonesResponse, err
		} else {
			transactionM = transaction
			if len(m.Recipients) > 0 {
				recipientsArray := []models.MilestonesRecipientResponse{}
				for _, r := range m.Recipients {
					user, _ := GetUserWithAccountID(extReq, r.AccountID)
					accountName := ""
					if user.ID != 0 {
						accountName = user.Lastname + " " + user.Firstname
					}
					recipientsArray = append(recipientsArray, models.MilestonesRecipientResponse{
						AccountID:   r.AccountID,
						AccountName: accountName,
						Email:       user.EmailAddress,
						PhoneNumber: user.PhoneNumber,
						Amount:      r.Amount,
					})
				}
				dDate, err := utility.FormatDate(transaction.DueDate, "2006-01-02", "2006-01-02 15::05")
				if err != nil {
					dDate = ""
				}
				milestonesResponse = append(milestonesResponse, models.MilestonesResponse{
					Index:            i + 1,
					MilestoneID:      milestoneID,
					Title:            m.Title,
					Amount:           m.Amount,
					Status:           m.Status,
					InspectionPeriod: strconv.Itoa(m.InspectionPeriod),
					DueDate:          dDate,
					Recipients:       recipientsArray,
				})
			}
		}

	}
	return transactionM, milestonesResponse, nil
}

func resolveCreateMilestoneTransaction(extReq request.ExternalRequest, milestones []models.MileStone, transactionAmount float64, transactionObj models.ResolveTransactionObj, db postgresql.Databases) (models.Transaction, []models.MilestonesResponse, error) {
	var (
		escrowCharge       = transactionAmount - getTotalAmoutForMilestones(milestones)
		milestonesResponse = []models.MilestonesResponse{}
		transactionM       = models.Transaction{}
	)

	for i, m := range milestones {
		milestoneID := utility.RandomString(20)
		transUssdCode := utility.GetRandomNumbersInRange(10000, 99999)
		dueDate, _ := utility.GetUnixString(m.DueDate, "2006-01-02", "2006-01-02")
		description := m.Description
		if description == "" {
			description = transactionObj.Description
		}
		quantity := m.Quantity
		if quantity == 0 {
			quantity = 1
		}
		gracePeriod, _ := utility.GetUnixString(m.GracePeriod, "2006-01-02", "2006-01-02")
		recipientsJson, err := json.Marshal(m.Recipients)
		if err != nil {
			return transactionM, milestonesResponse, err
		}

		transaction := models.Transaction{
			TransactionID:    transactionObj.TransactionID,
			PartiesID:        transactionObj.TransactionPartiesID,
			Title:            transactionObj.Title + ";" + m.Title + ";" + strconv.Itoa(int(transactionObj.Amount)) + ";" + strconv.Itoa(i+1),
			Type:             transactionObj.Type,
			Description:      description,
			MilestoneID:      milestoneID,
			Amount:           m.Amount,
			Status:           m.Status,
			Quantity:         quantity,
			InspectionPeriod: strconv.Itoa(m.InspectionPeriod),
			DueDate:          dueDate,
			ShippingFee:      m.ShippingFee,
			GracePeriod:      gracePeriod,
			Currency:         transactionObj.Currency,
			Country:          transactionObj.Country,
			BusinessID:       transactionObj.BusinessID,
			EscrowCharge:     escrowCharge,
			TransUssdCode:    transUssdCode,
			Recipients:       string(recipientsJson),
			DisputeHandler:   transactionObj.DisputeHandler,
			EscrowWallet:     transactionObj.EscrowWallet,
		}

		err = transaction.CreateTransaction(db.Transaction)
		if err != nil {
			return transactionM, milestonesResponse, err
		} else {
			transactionM = transaction
			if len(m.Recipients) > 0 {
				recipientsArray := []models.MilestonesRecipientResponse{}
				for _, r := range m.Recipients {
					user, _ := GetUserWithAccountID(extReq, r.AccountID)
					accountName := ""
					if user.ID != 0 {
						accountName = user.Lastname + " " + user.Firstname
					}
					recipientsArray = append(recipientsArray, models.MilestonesRecipientResponse{
						AccountID:   r.AccountID,
						AccountName: accountName,
						Email:       user.EmailAddress,
						PhoneNumber: user.PhoneNumber,
						Amount:      r.Amount,
					})
				}
				dDate, err := utility.FormatDate(transaction.DueDate, "2006-01-02", "2006-01-02 15::05")
				if err != nil {
					dDate = ""
				}
				milestonesResponse = append(milestonesResponse, models.MilestonesResponse{
					Index:            i + 1,
					MilestoneID:      milestoneID,
					Title:            m.Title,
					Amount:           m.Amount,
					Status:           m.Status,
					InspectionPeriod: strconv.Itoa(m.InspectionPeriod),
					DueDate:          dDate,
					Recipients:       recipientsArray,
				})
			}
		}

	}
	return transactionM, milestonesResponse, nil
}

func getTotalAmoutForMilestones(milestones []models.MileStone) float64 {
	var (
		totalAmount float64 = 0
	)

	for _, m := range milestones {
		totalAmount += m.Amount
	}
	return totalAmount
}

func resolveParties(extReq request.ExternalRequest, parties []models.Party, transactionPartiesID, transactionID string, db postgresql.Databases) ([]models.TransactionParty, []models.PartyResponse, error) {
	var (
		transactionParties = []models.TransactionParty{}
		partiesResponse    = []models.PartyResponse{}
	)

	for _, p := range parties {
		var roleCapabilities map[string]interface{}
		inrec, err := json.Marshal(&p.AccessLevel)
		if err != nil {
			return transactionParties, partiesResponse, err
		}
		err = json.Unmarshal(inrec, &roleCapabilities)
		if err != nil {
			return transactionParties, partiesResponse, err
		}
		user, _ := GetUserWithAccountID(extReq, p.AccountID)
		accountName := ""
		if user.ID != 0 {
			accountName = user.Lastname + " " + user.Firstname
		}

		transactionParty := models.TransactionParty{
			TransactionPartiesID: transactionPartiesID,
			TransactionID:        transactionID,
			AccountID:            p.AccountID,
			Role:                 p.Role,
			Status:               p.Status,
			RoleCapabilities:     roleCapabilities,
		}

		err = transactionParty.CreateTransactionParty(db.Transaction)
		if err != nil {
			return transactionParties, partiesResponse, err
		} else {
			transactionParties = append(transactionParties, transactionParty)
			partiesResponse = append(partiesResponse, models.PartyResponse{
				PartyID:     int(transactionParty.ID),
				AccountID:   transactionParty.AccountID,
				AccountName: accountName,
				PhoneNumber: user.PhoneNumber,
				Email:       user.EmailAddress,
				Role:        transactionParty.Role,
				Status:      transactionParty.Status,
				AccessLevel: p.AccessLevel,
			})
		}
	}

	return transactionParties, partiesResponse, nil
}

func resolveTransactionFiles(files []models.File, transactionID string, accountID int, db postgresql.Databases) ([]models.TransactionFile, error) {
	var (
		transactionFiles = []models.TransactionFile{}
	)
	for _, f := range files {
		transactionFile := models.TransactionFile{
			TransactionID: transactionID,
			AccountID:     accountID,
			FileUrl:       f.URL,
		}
		err := transactionFile.CreateTransactionFile(db.Transaction)
		if err != nil {
			return transactionFiles, err
		} else {
			transactionFiles = append(transactionFiles, transactionFile)
		}
	}

	return transactionFiles, nil
}

func GetTransactionStatus(index string) string {
	dataMap := map[string]string{
		"":        "Draft",
		"sac":     "Sent - Awaiting Confirmation",
		"sr":      "Sent - Rejected",
		"af":      "Accepted - Funded",
		"anf":     "Accepted - Not Funded",
		"fr":      "Funded - Rejected",
		"ip":      "In Progress",
		"d":       "Delivered",
		"da":      "Delivered - Accepted",
		"dr":      "Delivered - Rejected",
		"cdp":     "Closed - Disbursement Pending",
		"cmdp":    "Closed - Manual Disbursement Pending",
		"cdc":     "Closed - Disbursement Complete",
		"cd":      "Closed - Disputed",
		"cnf":     "Closed - Not Funded",
		"closed":  "Closed",
		"draft":   "Draft",
		"active":  "Active",
		"cr":      "Closed - Refunded",
		"deleted": "Deleted",
	}
	status := dataMap[strings.ToLower(index)]
	if status == "" {
		status = dataMap[""]
	}
	return status
}

func validatePartiesAndMilestones(tType string, parties []models.Party, milestones []models.MileStone) error {
	for _, v := range parties {
		if v.Role == "" {
			return fmt.Errorf("A role id is required for this transaction.")
		}
		if v.AccountID == 0 {
			return fmt.Errorf("An account id is required for %v", v.Role)
		}
		baseAccessLevel := models.PartyAccessLevel{}
		if v.AccessLevel == baseAccessLevel {
			return fmt.Errorf("An access level is required for %v", v.Role)
		}
	}

	if tType == "milestone" || tType == "oneoff" {
		if tType == "milestone" && len(milestones) < 2 {
			return fmt.Errorf("Your transaction type: %v needs at-least two milestone", tType)
		}
		if tType == "oneoff" && len(milestones) < 1 {
			return fmt.Errorf("Your transaction type: %v needs at-least one milestone", tType)
		}

		for i, v := range milestones {
			position := i + 1

			if v.Title == "" {
				return fmt.Errorf("Title is required for milestone %v", position)
			}
			if v.InspectionPeriod == 0 {
				return fmt.Errorf("Inspection is required for milestone %v", position)
			}
			if v.DueDate == "" {
				return fmt.Errorf("Due date is required for milestone %v", position)
			}
			_, err := validateDueDate(v.GracePeriod)
			if err != nil {
				return fmt.Errorf("incorrect grace period format, try 2006-01-15")
			}

			_, err = validateDueDate(v.DueDate)
			if err != nil {
				return fmt.Errorf("incorrect due date format, try 2006-01-15")
			}

			if len(v.Recipients) < 1 {
				return fmt.Errorf("Recipients is required for milestone %v", position)
			}

		}
	}

	return nil
}

func validateDueDate(dateString string) (string, error) {
	if dateString != "" {
		dateString, err := utility.FormatDate(dateString, "2006-01-02", "2006-01-02")
		if err != nil {
			return dateString, err
		}

	}
	return dateString, nil
}
