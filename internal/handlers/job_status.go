package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"retailpulse/internal/services"
	"strconv"
)

func JobStatusHandler(jobService *services.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get jobid from URL query
		jobIDStr := r.URL.Query().Get("jobid")

		validJobID, err := jobService.ValidateJobID(jobIDStr)
		if err != nil {
			http.Error(w," Unable to validate jobID ", http.StatusInternalServerError)
			return
		}
		if !validJobID {
			http.Error(w, "Invalid jobid ", http.StatusBadRequest)
			return
		}

		// Convert jobid to integer
		jobID, err := strconv.Atoi(jobIDStr)
		if err != nil {
			// If conversion fails, return an error
			http.Error(w, "Invalid jobid format", http.StatusBadRequest)
			return
		}
		// Get the job status
		response, err := jobService.GetJobStatus(jobID)
		if err != nil {
			// If job status retrieval fails, return an error
			http.Error(w, fmt.Sprintf("Failed to retrieve job status: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with the job status
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
