package paymentprocessor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	Pending  PaymentStatus = "pending"
	Success PaymentStatus = "success"
	Failure PaymentStatus = "failure"
)

// PaymentRequest represents a payment request
type PaymentRequest struct {
	ID        string      `json:"id"`
	Amount    float64     `json:"amount"`
	Currency  string      `json:"currency"`
	Card      CardDetails `json:"card"`
	Timestamp time.Time   `json:"timestamp"`
}

// CardDetails represents card details
type CardDetails struct {
	Number string `json:"number"`
	Expiry string `json:"expiry"`
	CVV    string `json:"cvv"`
}

// ValidatePaymentRequest validates a payment request
func ValidatePaymentRequest(request *PaymentRequest) error {
	if request == nil {
		return errors.New("payment request is nil")
	}
	if request.ID == "" {
		return errors.New("payment request id is empty")
	}
	if request.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}
	if request.Currency == "" {
		return errors.New("payment currency is empty")
	}
	if request.Card.Number == "" || request.Card.Expiry == "" || request.Card.CVV == "" {
		return errors.New("card details are incomplete")
	}
	return nil
}

// ProcessPayment processes a payment request
func ProcessPayment(request *PaymentRequest) (*PaymentStatus, error) {
	if err := ValidatePaymentRequest(request); err != nil {
		return nil, err
	}
	// Simulate payment processing
	time.Sleep(2 * time.Second)
	// Generate a random payment status
	status := Pending
	if request.Amount < 100 {
		status = Success
	} else {
		status = Failure
	}
	return &status, nil
}

// HandlePayment handles a payment request from an HTTP endpoint
func HandlePayment(w http.ResponseWriter, r *http.Request) {
	var request PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.ID == "" {
		request.ID = uuid.New().String()
	}
	request.Timestamp = time.Now()
	status, err := ProcessPayment(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": string(*status)}); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}