package config

import (
	"fmt"
	"log"
	"strings" // Package untuk manipulasi string

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	DBConn string `mapstructure:"DB_CONN"`
	Port   string `mapstructure:"PORT"`
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() (*Config, error) {
	// Set config file name and type
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")   // Look for config in the working directory
	viper.AddConfigPath("./")  // Also check current directory
	viper.AddConfigPath("../") // And parent directory

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_CONN", "")

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("⚠️  .env file not found, using environment variables and defaults")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		log.Printf("✅ Config loaded from: %s\n", viper.ConfigFileUsed())
	}

	// Create config struct
	config := &Config{
		DBConn: viper.GetString("DB_CONN"),
		Port:   viper.GetString("PORT"),
	}

	// Validate required fields
	if config.DBConn == "" {
		log.Println("⚠️  DB_CONN not set, will use default local PostgreSQL")
	}

	return config, nil
}

// GetDatabaseURL returns the database URL with fallback to default
func (c *Config) GetDatabaseURL() string {
	if c.DBConn != "" {
		connStr := c.DBConn

		// Jika sudah ada default_query_exec_mode, skip
		// Menggunakan strings.Contains() dari standard library
		if strings.Contains(connStr, "default_query_exec_mode") {
			log.Println("✅ Connection string already has default_query_exec_mode parameter")
			return connStr
		}

		// Tambahkan default_query_exec_mode=simple_protocol untuk disable prepared statements
		// Ini fix error "prepared statement already exists" yang sering muncul
		// dengan PostgreSQL connection pooler (Railway/Supabase/PgBouncer)
		separator := "?"
		if strings.Contains(connStr, "?") {
			// Sudah ada query parameters, tambahkan dengan &
			separator = "&"
		}

		connStr += separator + "default_query_exec_mode=simple_protocol"
		log.Println("✅ Added default_query_exec_mode=simple_protocol to connection string")

		return connStr
	}
	// Default local PostgreSQL connection
	return "host=localhost user=postgres password=postgres dbname=kasir_db port=5432 sslmode=disable default_query_exec_mode=simple_protocol"
}
