package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const baseURL = "http://localhost:8080"

var authToken string

func main() {
	fmt.Println("=== 🚀 Starting Feature Limits Acceptance Test ===")

	// 1. Login as default admin to get token
	loginReq := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", baseURL+"/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed to login: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dump, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Login failed (status %d): %s\n", resp.StatusCode, string(dump))
		os.Exit(1)
	}

	var loginResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)
	authToken = loginResp.Data.Token
	fmt.Println("✅ Logged in successfully.")

	// 2. Check limits via GET /api/store/limits
	fmt.Println("\n--- Test 1: Fetching current store limits ---")
	req, _ = http.NewRequest("GET", baseURL+"/api/store/limits", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
	resp, _ = client.Do(req)

	if resp.StatusCode != http.StatusOK {
		dump, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Failed to fetch limits (status %d): %s\n", resp.StatusCode, string(dump))
	} else {
		dump, _ := io.ReadAll(resp.Body)
		fmt.Printf("✅ Store limits fetched successfully: %s\n", string(dump))
	}
	resp.Body.Close()

	// 3. Try creating a cashier & admin limitations
	fmt.Println("\n--- Test 2: Creating Cashiers & Admin Limits ---")
	createUser := func(username, role string) (int, string) {
		cReq := map[string]string{
			"username": username,
			"password": "password123",
			"role":     role,
		}
		b, _ := json.Marshal(cReq)
		req, _ := http.NewRequest("POST", baseURL+"/api/users", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		resp, _ := client.Do(req)

		dump, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return resp.StatusCode, string(dump)
	}

	// Coba buat admin (harus gagal)
	adminStatus, adminDump := createUser("admin_limit_test", "admin")
	fmt.Printf("Create Admin admin_limit_test: Status %d - %s", adminStatus, adminDump)
	if adminStatus == http.StatusBadRequest {
		fmt.Println("✅ Admin creation successfully blocked!")
	} else {
		fmt.Println("⚠️ Did not receive 400 Bad Request when creating admin.")
	}

	// Assuming Free Package max_kasir = 1. We might already have kasirs.
	// We'll just try to create up to 3 kasirs and expect at least one 403 Forbidden.
	got403 := false
	for i := 1; i <= 3; i++ {
		status, dump := createUser(fmt.Sprintf("kasir_limit_test_%d", i), "kasir")
		fmt.Printf("Create Kasir %s: Status %d - %s", fmt.Sprintf("kasir_limit_test_%d", i), status, dump)
		if status == http.StatusForbidden {
			got403 = true
			if strings.Contains(dump, "Batas kasir untuk paket") {
				fmt.Println("✅ Kasir creation limit exact formatted string matching!")
			} else {
				fmt.Println("⚠️ 403 received but string formatting does not EXACTLY match 'Batas kasir untuk paket'")
			}
			fmt.Println("✅ Kasir creation successfully blocked by limits!")
			break
		}
	}
	if !got403 {
		fmt.Println("⚠️ Did not receive 403 Forbidden when creating kasirs. Make sure limits are enforcing.")
	}

	// 4. Test product creation limits via Purchases
	fmt.Println("\n--- Test 3: Creating Products via Purchase ---")

	// Create Category as prerequisite
	catReq := map[string]string{"nama": "Test Category Limit", "deskripsi": "Test"}
	bc, _ := json.Marshal(catReq)
	reqCat, _ := http.NewRequest("POST", baseURL+"/api/categories", bytes.NewBuffer(bc))
	reqCat.Header.Set("Content-Type", "application/json")
	reqCat.Header.Set("Authorization", "Bearer "+authToken)
	resCat, _ := client.Do(reqCat)
	resCatBody, _ := io.ReadAll(resCat.Body)
	resCat.Body.Close()

	var catData struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(resCatBody, &catData)
	catID := catData.Data.ID
	if catID == 0 {
		fmt.Printf("⚠️ Failed to create category, using default 1. Response: %s\n", string(resCatBody))
		catID = 1
	}

	createPurchaseWithNewProducts := func(numProducts int) (int, string, map[string]interface{}) {
		var items []map[string]interface{}
		for i := 0; i < numProducts; i++ {
			items = append(items, map[string]interface{}{
				"product_name": fmt.Sprintf("Limit Test Product %d", i+1),
				"quantity":     10,
				"buy_price":    1000,
				"sell_price":   1500,
				"category_id":  catID, // Use valid category
			})
		}

		pReq := map[string]interface{}{
			"supplier_id": nil,
			"items":       items,
		}
		b, _ := json.Marshal(pReq)
		req, _ := http.NewRequest("POST", baseURL+"/api/purchases", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		resp, _ := client.Do(req)

		dump, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var jsonRes map[string]interface{}
		json.Unmarshal(dump, &jsonRes)

		return resp.StatusCode, string(dump), jsonRes
	}

	// Paksa test limits dengan mencoba insert 11 produk baru sekaligus (Limit Gratis = 10)
	status, dump, _ := createPurchaseWithNewProducts(11)
	fmt.Printf("Create Purchase (adding 11 products): Status %d - %s\n", status, dump)
	if status == http.StatusForbidden {
		fmt.Println("✅ Product limits successfully blocked creating > max_products via purchases!")
	} else {
		fmt.Println("⚠️ Did not receive 403 Forbidden when creating too many products. Make sure limits are enforcing.")
	}

	// Fetch a valid product ID for checkout
	createPurchaseWithNewProducts(1) // Create one product to ensure at least one exists
	prodReq, _ := http.NewRequest("GET", baseURL+"/api/produk?limit=1", nil)
	prodReq.Header.Set("Authorization", "Bearer "+authToken)
	prodRes, _ := client.Do(prodReq)
	prodDump, _ := io.ReadAll(prodRes.Body)
	prodRes.Body.Close()

	var prodData struct {
		Data []struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(prodDump, &prodData)
	var validProductID float64 = 1
	if len(prodData.Data) > 0 {
		validProductID = prodData.Data[0].ID
		fmt.Printf("Fetched valid product ID: %.0f\n", validProductID)
	} else {
		fmt.Printf("⚠️ Failed to fetch valid product. Response: %s\n", string(prodDump))
	}

	// 5. Test checkout limits
	fmt.Println("\n--- Test 4: Checkout Limits ---")
	doCheckout := func() (int, string) {
		cReq := map[string]interface{}{
			"customer_id":    nil,
			"payment_amount": 50000,
			"items": []map[string]interface{}{
				{
					"product_id": int(validProductID),
					"quantity":   1,
					"price":      10000,
				},
			},
		}
		b, _ := json.Marshal(cReq)
		req, _ := http.NewRequest("POST", baseURL+"/api/checkout", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		resp, _ := client.Do(req)

		dump, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return resp.StatusCode, string(dump)
	}

	gotCheckoutBlock := false
	for i := 1; i <= 15; i++ {
		status, dump := doCheckout()
		if i == 1 {
			fmt.Printf("Checkout 1: Status %d - %s\n", status, dump)
		}
		if status == http.StatusForbidden {
			gotCheckoutBlock = true
			fmt.Printf("Checkout %d: Status %d - %s\n", i, status, dump)
			fmt.Println("✅ Checkout successfully blocked by daily sales limits!")
			break
		}
	}
	if !gotCheckoutBlock {
		fmt.Println("⚠️ Did not receive 403 Forbidden when spamming checkout. Make sure limits are enforcing.")
	}

	fmt.Println("\n=== 🎉 Acceptance Test Complete ===")
}
