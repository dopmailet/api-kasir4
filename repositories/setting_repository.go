package repositories

import (
	"database/sql"
	"encoding/json"
	"kasir-api/models"
	"log"
)

type SettingRepository struct {
	db *sql.DB
}

func NewSettingRepository(db *sql.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetCustomerSettings retrieves the combined settings for customer feature
func (r *SettingRepository) GetCustomerSettings() (*models.AppSettings, error) {
	var jsonValue string
	err := r.db.QueryRow("SELECT value_json FROM app_settings WHERE key = 'customer_settings'").Scan(&jsonValue)

	defaultSettings := &models.AppSettings{
		ShowCustomerInPOS:   true,
		EnableLoyaltyPoints: true,
	}

	if err == sql.ErrNoRows {
		// Return default, don't fail
		return defaultSettings, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonValue), defaultSettings)
	if err != nil {
		log.Printf("Bypass: Failed to unmarshal customer settings, using defaults. Err: %v\n", err)
	}

	return defaultSettings, nil
}

// UpdateCustomerSettings updates or inserts the settings
func (r *SettingRepository) UpdateCustomerSettings(settings *models.AppSettings) error {
	jsonBytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO app_settings (key, value_json, updated_at) 
		VALUES ('customer_settings', $1, NOW())
		ON CONFLICT (key) DO UPDATE 
		SET value_json = $1, updated_at = NOW()
	`
	_, err = r.db.Exec(query, string(jsonBytes))
	return err
}
