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

// GetCustomerSettings retrieves the combined settings for customer feature for a specific store
func (r *SettingRepository) GetCustomerSettings(storeID int) (*models.AppSettings, error) {
	var jsonValue string
	err := r.db.QueryRow("SELECT value_json FROM app_settings WHERE key = 'customer_settings' AND store_id = $1", storeID).Scan(&jsonValue)

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

// UpdateCustomerSettings updates or inserts the settings for a specific store
func (r *SettingRepository) UpdateCustomerSettings(storeID int, settings *models.AppSettings) error {
	jsonBytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO app_settings (store_id, key, value_json, updated_at) 
		VALUES ($1, 'customer_settings', $2, NOW())
		ON CONFLICT (store_id, key) DO UPDATE 
		SET value_json = EXCLUDED.value_json, updated_at = NOW()
	`
	_, err = r.db.Exec(query, storeID, string(jsonBytes))
	return err
}

// GetPlatformSettings retrieves global platform settings (used by superadmin)
func (r *SettingRepository) GetPlatformSettings() (*models.PlatformSettings, error) {
	var jsonValue string
	// Kita simpan di store_id = 1 agar cocok dengan constraint tabel multi-tenant saat ini
	err := r.db.QueryRow("SELECT value_json FROM app_settings WHERE key = 'platform_settings' AND store_id = 1").Scan(&jsonValue)

	defaultSettings := &models.PlatformSettings{
		AdminWhatsApp: nil,
	}

	if err == sql.ErrNoRows {
		// Return default, don't fail
		return defaultSettings, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonValue), defaultSettings)
	if err != nil {
		log.Printf("Bypass: Failed to unmarshal platform settings, using defaults. Err: %v\n", err)
	}

	return defaultSettings, nil
}

// UpdatePlatformSettings updates or inserts the global platform settings
func (r *SettingRepository) UpdatePlatformSettings(settings *models.PlatformSettings) error {
	jsonBytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO app_settings (store_id, key, value_json, updated_at) 
		VALUES (1, 'platform_settings', $1, NOW())
		ON CONFLICT (store_id, key) DO UPDATE 
		SET value_json = EXCLUDED.value_json, updated_at = NOW()
	`
	_, err = r.db.Exec(query, string(jsonBytes))
	return err
}
