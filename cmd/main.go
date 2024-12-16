package main

import (
    "log"
    "net/http"
    "retailpulse/internal/handlers"
    "retailpulse/internal/services"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "os"
    "fmt"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get DB credentials from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    portNumber := os.Getenv("PORT_NUMBER")

    // Construct the database connection string
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

    // Open DB connection
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Initialize services with DB connection
    jobService := services.NewJobService(db)

    // Set up HTTP handlers
    http.HandleFunc("/api/submit/", handlers.SubmitJobHandler(jobService))
    http.HandleFunc("/api/status", handlers.JobStatusHandler(jobService))

    // Start the server
    log.Printf("Server is listening on port %s", portNumber)
    log.Fatal(http.ListenAndServe(":" + portNumber, nil))
}
