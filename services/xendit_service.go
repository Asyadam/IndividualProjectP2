package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type XenditService interface {
	CreateInvoice(externalID string, amount int, payerEmail string, description string) (string, error)
}

type xenditService struct{}

func NewXenditService() XenditService {
	return &xenditService{}
}

type xenditInvoiceRequest struct {
	ExternalID  string `json:"external_id"`
	Amount      int    `json:"amount"`
	PayerEmail  string `json:"payer_email"`
	Description string `json:"description"`
}

type xenditInvoiceResponse struct {
	ID         string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
	Status     string `json:"status"`
}

func (s *xenditService) CreateInvoice(externalID string, amount int, payerEmail string, description string) (string, error) {
	baseURL := os.Getenv("XENDIT_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.xendit.co"
	}

	secretKey := os.Getenv("XENDIT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("xendit secret key is empty")
	}

	payload := xenditInvoiceRequest{
		ExternalID:  externalID,
		Amount:      amount,
		PayerEmail:  payerEmail,
		Description: description,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v2/invoices", baseURL),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(secretKey, "")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("failed to create xendit invoice")
	}

	var invoiceResponse xenditInvoiceResponse

	if err := json.NewDecoder(resp.Body).Decode(&invoiceResponse); err != nil {
		return "", err
	}

	if invoiceResponse.InvoiceURL == "" {
		return "", errors.New("xendit invoice url is empty")
	}

	return invoiceResponse.InvoiceURL, nil
}
