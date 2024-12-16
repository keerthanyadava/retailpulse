package utils

import (
	"bytes"
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg" // Supports JPEG format
	_ "image/png"  // Supports PNG format
	"io"
	"math/rand"
	"net/http"
	"time"
)

// ImageProcessor calculates the perimeter of an image and simulates GPU processing with a random sleep
func ImageProcessor(imageURL string) (float64, error) {
	// Download the image
	resp, err := http.Get(imageURL)
	if err != nil {
		return 0, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read image data: %v", err)
	}

	// Decode the image data
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %v", err)
	}

	// Calculate the perimeter
	bounds := img.Bounds()
	width, height := float64(bounds.Max.X-bounds.Min.X), float64(bounds.Max.Y-bounds.Min.Y)
	perimeter := 2 * (width + height)

	// Simulate GPU processing with a random sleep
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sleepDuration := time.Duration(r.Intn(300)+100) * time.Millisecond
	time.Sleep(sleepDuration)

	return perimeter, nil
}

// StoreImageResult inserts a processed image result into the database
func StoreImageResult(db *sql.DB, jobID int, imageURL string, perimeter float64, storeID string) error {
	// Prepare the SQL statement
	query := "INSERT INTO images (job_id, url, perimeter, store_id) VALUES (?, ?, ?, ?)"

	// Execute the SQL query to store the result
	_, err := db.Exec(query, jobID, imageURL, perimeter, storeID)
	if err != nil {
		return fmt.Errorf("failed to store image result: %v", err)
	}

	return nil
}
