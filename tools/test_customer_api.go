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
	doRequest := func(method, path string, body interface{}) (interface{}, int) {
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
			return nil, 0
		}
		defer resp.Body.Close()

		var resBody interface{}
		json.NewDecoder(resp.Body).Decode(&resBody)
		return resBody, resp.StatusCode
	}

	// 2. Create Customer
	log.Println("\n--- 2. Create Customer ---")
	custPayload := map[string]interface{}{
		"name":    "Budi Test",
		"phone":   "081234567890",
		"address": "Jalan Test 123",
	}
	res, code := doRequest("POST", "/customers", custPayload)
	var custID float64
	if code == 201 {
		log.Println("✅ Create Customer Berhasil")
		custData := res.(map[string]interface{})
		custID = custData["id"].(float64)
	} else {
		log.Println("❌ Create Customer Gagal:", res)
	}

	// 3. Search Customer
	log.Println("\n--- 3. Search Customer ---")
	res, code = doRequest("GET", "/customers/search?q=Budi", nil)
	if code == 200 {
		log.Println("✅ Search Customer Berhasil")
	} else {
		log.Println("❌ Search Customer Gagal:", res)
	}

	// 4. Test Checkout with Customer
	if custID > 0 {
		log.Println("\n--- 4. Checkout with Customer (ID:", int(custID), ") ---")
		checkoutPayload := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"product_id":      1,
					"quantity":        2,
					"price":           10000,
					"discount_type":   "",
					"discount_value":  0,
					"discount_amount": 0,
				},
			},
			"payment_amount":  50000,
			"discount_amount": 0,
			"customer_id":     int(custID),
		}
		res, code = doRequest("POST", "/checkout", checkoutPayload)
		if code == 200 || code == 201 {
			log.Println("✅ Checkout dengan Customer Berhasil")
			resMap, _ := res.(map[string]interface{})
			log.Printf("Response points: %+v\n", resMap["data"])
		} else {
			log.Println("❌ Checkout Gagal:", res)
		}

		// 5. Cek history loyalty
		log.Println("\n--- 5. Get Loyalty History ---")
		res, code = doRequest("GET", fmt.Sprintf("/customers/%d/loyalty-transactions", int(custID)), nil)
		if code == 200 {
			log.Println("✅ Get Loyalty history Berhasil")
			log.Printf("Data: %+v\n", res)
		} else {
			log.Println("❌ Get Loyalty history Gagal:", res)
		}
	}

	log.Println("\nSelesai Pengujian Dasar!")
}
