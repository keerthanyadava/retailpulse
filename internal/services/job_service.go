package services

import (
	"database/sql"
	"fmt"
	"sync"

	"retailpulse/internal/models"
	"retailpulse/internal/utils"
)

type JobService struct {
	db *sql.DB
}

func NewJobService(db *sql.DB) *JobService {
	return &JobService{
		db: db,
	}
}

func (js *JobService) SubmitJob(request models.SubmitJobRequest) (int, error) {
	// Validate job request

	if request.Count != len(request.Visits) {
		return 0, fmt.Errorf("count mismatch: expected %d visits, but got %d", request.Count, len(request.Visits))
	}

	// Ensure each visit has a valid store_id and image_url
	for _, visit := range request.Visits {
		// Check if store_id is present
		if visit.StoreID == "" {
			return 0, fmt.Errorf("missing store_id for visit")
		}
		// Check if the store exists (assuming a function exists to validate the store)
		storeExists, err := utils.CheckStoreExists(js.db, visit.StoreID)
		if err != nil {
			return 0, fmt.Errorf("failed to check store existence: %v", err)
		}
		if !storeExists {
			return 0, fmt.Errorf("invalid store_id: %s ,store does not exist", visit.StoreID)
		}
		// Validate image_url is not empty
		if len(visit.ImageURLs) == 0 {
			return 0, fmt.Errorf("missing image_url for store %s", visit.StoreID)
		}
	}

	// Insert initial job status using utility function

	jobID, err := utils.RegisterJob(js.db)
	if err != nil {
		return 0, fmt.Errorf("failed to register job: %v", err)
	}

	// Process the job asynchronously
	go js.processJob(jobID, request)

	return jobID, nil
}

func (js *JobService) processJob(jobID int, request models.SubmitJobRequest) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	finalStatus := "completed"

	for _, visit := range request.Visits {
		wg.Add(1)
		go func(visit models.Visit) {
			defer wg.Done()
			storeError, ok := processVisit(js.db, visit, jobID)

			if !ok {

				mu.Lock()
				utils.LogJobFailure(js.db, jobID, storeError)
				finalStatus = "failed"
				mu.Unlock()
			}
		}(visit)
	}
	wg.Wait()

	err := utils.UpdateJobStatus(js.db, jobID, finalStatus)
	if err != nil {
		fmt.Printf("Failed to update job status for job %d: %v\n", jobID, err)
	}

}

func (js *JobService) GetJobStatus(jobID int) (*models.JobStatusResponse, error) {
	// Call the utility function to fetch the job status and error details
	jobData, err := utils.GetJobStatus(js.db, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch job status: %v", err)
	}

	// Extract status
	status, ok := jobData["status"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid status format")
	}

	// Convert failures to the correct type
	var errors []models.StoreError
	if status == "failed" {
		failureData, ok := jobData["error"].([]map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid error format")
		}

		for _, failure := range failureData {

			storeID, _ := failure["store_id"].(string)
			errorMsg, _ := failure["error"].(string)

			errors = append(errors, models.StoreError{
				StoreID: storeID,
				Error:   errorMsg,
			})
		}
	}

	// Ensure job ID is correct
	retrievedJobID, ok := jobData["job_id"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid job ID format")
	}

	// Construct the JobStatusResponse
	response := &models.JobStatusResponse{
		Status: status,
		JobID:  retrievedJobID,
		Errors: errors,
	}

	return response, nil
}

func processVisit(db *sql.DB, visit models.Visit, jobID int) (models.StoreError, bool) {
	var storeErrors models.StoreError
	// Check store existence from the database
	_, err := utils.FindStoreByID(db, visit.StoreID)
	if err != nil {
		storeErrors = models.StoreError{
			StoreID: visit.StoreID,
			Error:   "store not found",
		}
		return storeErrors, false
	}
	// Process images using the new image processor function
	for _, imageURL := range visit.ImageURLs {
		err := processImage(db, visit.StoreID, imageURL, jobID)
		if err != nil {
			storeErrors = models.StoreError{
				StoreID: visit.StoreID,
				Error:   err.Error(),
			}
			return storeErrors, false
		}
	}
	return storeErrors, true
}

func processImage(db *sql.DB, storeID, imageURL string, jobID int) error {
	// Call the ImageProcessor function to handle image downloading, perimeter calculation, and storage

	perimeter, err := utils.ImageProcessor(imageURL)
	if err != nil {
		return fmt.Errorf("image processing failed for image %s: %v", imageURL, err)
	}

	// Store the processed result (image perimeter)
	err = utils.StoreImageResult(db, jobID, imageURL, perimeter, storeID)
	if err != nil {
		return fmt.Errorf("failed to store image result for image %s: %v", imageURL, err)
	}

	return nil
}


func (js *JobService) ValidateJobID(jobID string) (bool, error) {
	valid, err:= utils.CheckJobExists(js.db, jobID)
	if err != nil {
		return valid, err
	}
	return valid, err
}