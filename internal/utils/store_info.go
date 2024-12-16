package utils

import (
    "database/sql"
    "errors"
    "retailpulse/internal/models"
    "fmt"
)

// FindStoreByID now queries the database to find the store by its ID
func FindStoreByID(db *sql.DB, storeID string) (*models.Store, error) {
    var store models.Store

    query := "SELECT StoreID, StoreName, AreaCode FROM stores WHERE StoreID = ?"
    row := db.QueryRow(query, storeID)
    
    err := row.Scan(&store.StoreID, &store.StoreName, &store.AreaCode)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("store not found")
        }
        return nil, fmt.Errorf("database query error: %v", err)
    }

    return &store, nil
}

// checkStoreExists checks if the store with the given store_id exists in the database
func CheckStoreExists(db *sql.DB ,storeID string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM stores WHERE StoreID = ?", storeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database query error: %v", err)
	}
	return count > 0, nil
}

// CheckJobExists checks if the job with the given jobID exists in the database
func CheckJobExists(db *sql.DB ,jobID string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM jobs WHERE id = ?", jobID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database query error: %v", err)
	}
	return count > 0, nil
}