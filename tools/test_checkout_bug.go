package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	baseURL := "http://localhost:8080/api"

	// 1. Get Token
	log.Println("--- 1. Login Admin ---")
	loginPayload := map[string]string{"username": "admin", "password": "admin123"}
	loginBody, _ := json.Marshal(loginPayload)
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		log.Fatal("Gagal login:", err)
	}
	defer resp.Body.Close()

	var loginRes map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginRes)

	if data, ok := loginRes["data"].(map[string]interface{}); ok {
		token := data["token"].(string)
		os.Setenv("TEST_TOKEN", token)
		log.Println("✅ Token didapat")
	} else {
		log.Fatal("❌ Gagal dapat token", loginRes)
	}

	token := os.Getenv("TEST_TOKEN")

	// Helper func
	doRequest := func(method, path string, body interface{}) (map[string]interface{}, int, string) {
		var reqBody io.Reader
		if body != nil {
			b, _ := json.Marshal(body)
			reqBody = bytes.NewBuffer(b)
		}
		req, _ := http.NewRequest(method, baseURL+path, reqBody)
		req.Header.Set("Authorization", "Bearer "+token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, 0, ""
		}
		defer resp.Body.Close()

		rawBody, _ := io.ReadAll(resp.Body)
		var jsonBody map[string]interface{}
		json.Unmarshal(rawBody, &jsonBody)

		return jsonBody, resp.StatusCode, string(rawBody)
	}

	customerId := 5 // Assuming customer 5 exists based on yesterday's run

	// Checkout
	log.Println("\n--- Checkout ---")
	checkoutPayload := map[string]interface{}{
		"customer_id":    customerId,
		"payment_amount": 100000,
		"items": []map[string]interface{}{
			{
				"product_id": 1,
				"quantity":   1,
				"price":      25000, // Make sure finalTotal > 10000
			},
		},
	}

	res, code, raw := doRequest("POST", "/checkout", checkoutPayload)
	log.Printf("Code: %d", code)
	log.Printf("Res: %v", res)
	log.Printf("Raw: %s", raw)

	// Get Customer
	log.Printf("\n--- Get Customer %d ---", customerId)
	res2, code2, _ := doRequest("GET", fmt.Sprintf("/customers/%d", customerId), nil)
	log.Printf("Code: %d", code2)
	log.Printf("Customer Data: %v", res2)
}
