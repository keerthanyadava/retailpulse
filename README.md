Retail Pulse Image Processing Service

Description
A Go-based microservice for processing store images, calculating image perimeters, and managing job statuses asynchronously. The service simulates job execution with random processing times and tracks job status in memory.

Key Features
1. Asynchronous Job Submission and Processing
2. Image URL Validation and Download
3. Random Processing Time Simulation
4. Concurrent Job Handling

Assumptions
- Basic error handling for job submission, image processing, and status updates.
- Limited error handling for store and image processing.

Setup & Installation

Prerequisites
- Go 1.23+ (for local setup)

Local Setup
1. Download Go modules:
   go mod download
2. Build the application:
   go build
3. Run :
    go run cmd/main.go
4. Setup Mysql DB using the db_schema
   The service will start on localhost:8080.

Improvements with More Time
1. Persistent Job Storage using CRUD interface
2. Enhanced Error Handling
3. Better Logging
4. Rate Limiting
5. Distributed Job Processing
6. Security Enhancements
