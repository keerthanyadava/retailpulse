package models

type SubmitJobRequest struct {
    Count   int     `json:"count"`
    Visits  []Visit `json:"visits"`
}

type Visit struct {
    StoreID    string   `json:"store_id"`
    ImageURLs  []string `json:"image_url"`
    VisitTime  string   `json:"visit_time"`
}

type JobResponse struct {
    JobID int `json:"job_id"`
}

type JobStatusResponse struct {
    Status string        `json:"status"`
    JobID  int        `json:"job_id"`
    Errors []StoreError  `json:"error,omitempty"`
}

type StoreError struct {
    StoreID string `json:"store_id"`
    Error   string `json:"error"`
}