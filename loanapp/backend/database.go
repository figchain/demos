package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection and creates tables
func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./loanapp.db")
	if err != nil {
		return err
	}

	// Create table
	createTableSQL := `CREATE TABLE IF NOT EXISTS loan_applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		applicant_name TEXT NOT NULL,
		credit_score INTEGER NOT NULL,
		amount_requested REAL NOT NULL,
		risk_factor REAL NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// GetAllApplications retrieves all loan applications from the database
func GetAllApplications() ([]LoanApplication, error) {
	rows, err := db.Query(`
		SELECT id, applicant_name, credit_score, amount_requested, risk_factor, status, created_at
		FROM loan_applications
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []LoanApplication
	for rows.Next() {
		var app LoanApplication
		err := rows.Scan(
			&app.ID,
			&app.ApplicantName,
			&app.CreditScore,
			&app.AmountRequested,
			&app.RiskFactor,
			&app.Status,
			&app.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, nil
}

// InsertApplication inserts a new loan application into the database
func InsertApplication(app LoanApplication) error {
	_, err := db.Exec(`
		INSERT INTO loan_applications (applicant_name, credit_score, amount_requested, risk_factor, status)
		VALUES (?, ?, ?, ?, ?)
	`, app.ApplicantName, app.CreditScore, app.AmountRequested, app.RiskFactor, app.Status)
	return err
}

// UpdateApplicationStatus updates the status of an application
func UpdateApplicationStatus(id int, status string) error {
	_, err := db.Exec(`
		UPDATE loan_applications
		SET status = ?
		WHERE id = ?
	`, status, id)
	return err
}
