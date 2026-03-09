package main

import (
	"bytes"
	"encoding/json"
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
	doRequest := func(method, path string, body interface{}) (interface{}, int, string) {
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
		var jsonBody interface{}
		json.Unmarshal(rawBody, &jsonBody)

		return jsonBody, resp.StatusCode, string(rawBody)
	}

	log.Println("\n--- Test 1. GET /api/customers/search ---")
	res, code, raw := doRequest("GET", "/customers/search?q=a", nil)
	log.Printf("Code: %d", code)
	log.Printf("Raw: %s", raw)

	log.Println("\n--- Test 2. GET /api/settings/customer ---")
	res, code, raw = doRequest("GET", "/settings/customer", nil)
	log.Printf("Code: %d", code)
	log.Printf("Raw: %s", raw)
	_ = res
}
