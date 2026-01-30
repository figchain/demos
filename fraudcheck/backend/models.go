package main

import (
	"time"
)

// LoanApplication represents a loan application record
type LoanApplication struct {
	ID              int       `json:"id"`
	ApplicantName   string    `json:"applicant_name"`
	CreditScore     int       `json:"credit_score"`
	AmountRequested float64   `json:"amount_requested"`
	RiskFactor      float64   `json:"risk_factor"` // 0.0 to 1.0
	Status          string    `json:"status"`      // approved, review_required, rejected
	CreatedAt       time.Time `json:"created_at"`
}
