package dto

type DepositRequest struct {
	Amount int `json:"amount" validate:"required,gt=0"`
}

type XenditCallbackRequest struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Amount     int    `json:"amount"`
}
