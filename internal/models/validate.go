package models

type ValidateOnDBReq struct {
	Table string      `validate:"required" json:"table"`
	Type  string      `validate:"required" json:"type"`
	Query string      `validate:"required" json:"query"`
	Value interface{} `json:"value"`
}
