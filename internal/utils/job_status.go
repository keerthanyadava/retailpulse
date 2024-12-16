package utils

import (
	"database/sql"
	"fmt"
	"log"
	"retailpulse/internal/models"
)

// RegisterJob inserts a new job with the 'ongoing' status and returns the job ID
func RegisterJob(db *sql.DB) (int, error) {

	// Insert a new job with the status 'ongoing'
	query := "INSERT INTO jobs (status) VALUES ('ongoing')"
	result, err := db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("could not insert job: %v", err)
	}

	// Get the auto-generated job ID
	jobID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not retrieve last inserted job ID: %v", err)
	}

	return int(jobID), nil
}

// UpdateJobStatus updates the job status in the jobs table
func UpdateJobStatus(db *sql.DB, jobID int, status string) error {
	// Update job status directly in the jobs table

	_, err := db.Exec(`
		UPDATE jobs 
		SET status = ? 
		WHERE id = ?
	`, status, jobID)
	if err != nil {
		return fmt.Errorf("could not update job status: %v", err)
	}

	return nil
}

// GetJobStatus retrieves the job status from the jobs table
func GetJobStatus(db *sql.DB, jobID int) (map[string]interface{}, error) {
	// Fetch the job status
	var jobStatus string
	query := "SELECT status FROM jobs WHERE id = ?"
	err := db.QueryRow(query, jobID).Scan(&jobStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch job status: %v", err)
	}

	// If the job status is "failed", fetch all associated failures
	if jobStatus == "failed" {
		failuresQuery := "SELECT store_id, error_message FROM job_failures WHERE job_id = ?"
		rows, err := db.Query(failuresQuery, jobID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch job failures: %v", err)
		}
		defer rows.Close()

		var failures []map[string]interface{}
		for rows.Next() {
			var storeID, errorMessage string
			if err := rows.Scan(&storeID, &errorMessage); err != nil {
				return nil, fmt.Errorf("failed to scan job failure row: %v", err)
			}

			failures = append(failures, map[string]interface{}{
				"store_id": storeID,
				"error":    errorMessage,
			})
		}

		return map[string]interface{}{
			"status": "failed",
			"job_id": jobID,
			"error":  failures,
		}, nil
	}

	// Handle other statuses (e.g., "success")
	return map[string]interface{}{
		"status": "success",
		"job_id": jobID,
	}, nil
}

func LogJobFailure(db *sql.DB, jobID int, storeError models.StoreError) error {
	query := "INSERT INTO job_failures (job_id, store_id, error_message) VALUES (?, ?, ?)"
	_, err := db.Exec(query, jobID, storeError.StoreID, storeError.Error)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("failed to log job failure: %v", err)
	}
	return nil
}
