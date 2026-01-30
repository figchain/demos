package main

import (
	"log"
	"math/rand"
	"time"
)

var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Barbara", "David", "Elizabeth", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Charles", "Karen", "Christopher", "Nancy", "Daniel", "Lisa",
	"Matthew", "Betty", "Anthony", "Margaret", "Mark", "Sandra", "Donald", "Ashley",
	"Steven", "Kimberly", "Paul", "Emily", "Andrew", "Donna", "Joshua", "Michelle",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
	"Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White",
	"Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young",
	"Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores",
}

// SeedDatabase populates the database with fake loan application data
func SeedDatabase() error {
	log.Println("Seeding database with fake loan applications...")

	// Check if data already exists
	apps, err := GetAllApplications()
	if err != nil {
		return err
	}

	if len(apps) > 0 {
		log.Printf("Database already contains %d applications, skipping seed", len(apps))
		return nil
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate 50 fake loan applications
	for i := 0; i < 50; i++ {
		app := generateFakeApplication()
		err := InsertApplication(app)
		if err != nil {
			return err
		}
	}

	log.Println("Successfully seeded 50 loan applications")
	return nil
}

// generateFakeApplication creates a random loan application
func generateFakeApplication() LoanApplication {
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]

	creditScore := rand.Intn(551) + 300                 // 300-850
	amountRequested := float64(rand.Intn(95000) + 5000) // $5,000 - $100,000

	// Calculate risk factor based on credit score and amount
	riskFactor := calculateRiskFactor(creditScore, amountRequested)

	// Determine status based on risk factor and current threshold
	threshold := float64(GetFailThreshold())
	riskPercentage := riskFactor * 100
	status := determineStatusByThreshold(riskPercentage, threshold)

	return LoanApplication{
		ApplicantName:   firstName + " " + lastName,
		CreditScore:     creditScore,
		AmountRequested: amountRequested,
		RiskFactor:      riskFactor,
		Status:          status,
	}
}

// calculateRiskFactor determines risk based on credit score and amount
func calculateRiskFactor(creditScore int, amount float64) float64 {
	// Base risk from credit score (inverted - lower score = higher risk)
	creditRisk := 1.0 - (float64(creditScore-300) / 550.0)

	// Additional risk from loan amount (higher amount = higher risk)
	amountRisk := amount / 100000.0

	// Combine risks (weighted average)
	totalRisk := (creditRisk * 0.7) + (amountRisk * 0.3)

	// Add some randomness (+/- 0.1)
	randomFactor := (rand.Float64() * 0.2) - 0.1
	totalRisk += randomFactor

	// Ensure risk is between 0 and 1
	if totalRisk < 0 {
		totalRisk = 0
	}
	if totalRisk > 1 {
		totalRisk = 1
	}

	return totalRisk
}

// determineStatus assigns a status based on risk factor
func determineStatus(riskFactor float64) string {
	if riskFactor < 0.3 {
		return "approved"
	} else if riskFactor < 0.7 {
		return "review_required"
	}
	return "rejected"
}
