package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"retailpulse/internal/models"
	"retailpulse/internal/services"
)

func SubmitJobHandler(jobService *services.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var request models.SubmitJobRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jobID, err := jobService.SubmitJob(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := models.JobResponse{JobID: jobID}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
