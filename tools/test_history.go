package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Ganti sesuai URL backend Anda
// const baseURL = "http://localhost:8080" // Local
const baseURL = "https://kasir-api-production-zainalzxc.up.railway.app" // Production

func main() {
	fmt.Println("🚀 KASIR API HISTORY TEST")
	fmt.Println("===================================")
	fmt.Println("Backend URL:", baseURL)

	// 1. Login Kasir
	kasirToken := login("kasir1", "kasir123")
	if kasirToken == "" {
		fmt.Println("❌ Failed to login kasir. Check Supabase credentials.")
		return
	}
	fmt.Println("✅ Kasir logged in")

	// 2. Fetch History
	fetchHistory(kasirToken)
}

func login(username, password string) string {
	fmt.Printf("logging in as %s...\n", username)
	body := map[string]string{"username": username, "password": password}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error login:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Login failed, status:", resp.StatusCode)
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response:", string(respBody))
		return ""
	}

	var result struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Data.Token
}

func fetchHistory(token string) {
	fmt.Println("Fetching transaction history...")

	req, _ := http.NewRequest("GET", baseURL+"/api/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error fetching history:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success! HTTP 200 OK")
		fmt.Println("Raw JSON Response:")
		fmt.Println(string(body))

		// Parse Array of Transactions
		var transactions []map[string]interface{}
		if err := json.Unmarshal(body, &transactions); err == nil {
			fmt.Println("-----------------------------------")
			fmt.Printf("Found %d transactions:\n", len(transactions))
			for i, t := range transactions {
				if i < 5 { // Show top 5 only
					fmt.Printf("- ID: %.0f | Total: %.2f | Date: %s\n", t["id"], t["total_amount"], t["created_at"])
				}
			}
			if len(transactions) > 5 {
				fmt.Println("... and more")
			}
		}
	} else {
		fmt.Println("❌ Failed to fetch history. Status:", resp.StatusCode)
		fmt.Println("Response:", string(body))
	}
}
