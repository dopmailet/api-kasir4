package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("🚀 KASIR API CATEGORY DISCOUNT TEST")
	fmt.Println("===================================")

	// 1. Login Admin
	adminToken := login("admin", "admin123")
	if adminToken == "" {
		fmt.Println("❌ Failed to login admin")
		return
	}
	fmt.Println("✅ Admin logged in")

	// 2. Create NEW Product for Test (ID auto-generated, assumed last inserted)
	// We need to fetch it or assume ID. Let's create a product with unique name.
	productID := createTestProduct(adminToken)
	fmt.Println("✅ Created Test Product ID:", productID)

	// 3. Create Discount for CATEGORY ID 1 (50%)
	createCategoryDiscount(adminToken, 1, 50.0)

	// 3. Login Kasir
	kasirToken := login("kasir1", "kasir123")
	if kasirToken == "" {
		fmt.Println("❌ Failed to login kasir")
		return
	}
	fmt.Println("✅ Kasir logged in")

	// 4. Checkout Product (Qty 2)
	// Price 10,000 * 2 = 20,000. Discount 50% = 10,000. Final = 10,000.
	checkout(kasirToken, productID, 2)
}

func createTestProduct(token string) int {
	productName := fmt.Sprintf("Test Product %d", time.Now().Unix())
	product := map[string]interface{}{
		"nama":        productName,
		"harga":       10000,
		"stok":        100,
		"category_id": 1, // Category 1
	}
	jsonBody, _ := json.Marshal(product)

	req, _ := http.NewRequest("POST", baseURL+"/api/produk", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating product:", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		fmt.Println("Failed to create product:", resp.StatusCode)
		return 0
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	// Note: API Create Product might return structure differently.
	// Let's assume standard response or parse carefully.
	if err := json.Unmarshal(body, &result); err != nil || result.Data.ID == 0 {
		// Fallback: parsing different structure or query by name?
		// For simplicity, let's assume we can get ID from response.
		// If API returns flat JSON:
		var flat struct {
			ID int `json:"id"`
		}
		json.Unmarshal(body, &flat)
		return flat.ID
	}
	return result.Data.ID
}

func createCategoryDiscount(token string, categoryID int, value float64) {
	fmt.Printf("Generating Discount for Category %d (%.0f%%)...\n", categoryID, value)

	discount := map[string]interface{}{
		"name":             fmt.Sprintf("Auto Test Discount Category %d", categoryID),
		"type":             "PERCENTAGE",
		"value":            value,
		"min_order_amount": 0,
		"category_id":      categoryID, // Category ID Logic
		"product_id":       nil,        // Ensure Product ID is nil
		"start_date":       time.Now().Format(time.RFC3339),
		"end_date":         time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"is_active":        true,
	}
	jsonBody, _ := json.Marshal(discount)

	req, _ := http.NewRequest("POST", baseURL+"/api/discounts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating discount:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 201 {
		fmt.Println("✅ Category Discount Created Successfully!")
	} else {
		fmt.Println("⚠️ Create Discount Response:", resp.StatusCode, string(body))
	}
}

func login(username, password string) string {
	body := map[string]string{"username": username, "password": password}
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error login:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Login failed, status:", resp.StatusCode)
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

func createDiscount(token string, productID int, value float64) {
	fmt.Printf("Generating Discount for Product %d (%.0f%%)...\n", productID, value)

	discount := map[string]interface{}{
		"name":             fmt.Sprintf("Auto Test Discount Product %d", productID),
		"type":             "PERCENTAGE",
		"value":            value,
		"min_order_amount": 0,
		"product_id":       productID,
		"start_date":       time.Now().Format(time.RFC3339),
		"end_date":         time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"is_active":        true,
	}
	jsonBody, _ := json.Marshal(discount)

	req, _ := http.NewRequest("POST", baseURL+"/api/discounts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating discount:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 201 {
		fmt.Println("✅ Discount Created Successfully!")
	} else {
		// Ignore error if it's likely a duplicate/already exists, just warn
		fmt.Println("⚠️ Create Discount Response:", resp.StatusCode, string(body))
	}
}

func checkout(token string, productID, qty int) {
	fmt.Printf("Checking out Product %d (Qty: %d)...\n", productID, qty)

	checkoutReq := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": productID, "quantity": qty},
		},
	}
	jsonBody, _ := json.Marshal(checkoutReq)

	req, _ := http.NewRequest("POST", baseURL+"/api/checkout", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error checkout:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Println("✅ Checkout Success! (Status:", resp.StatusCode, ")")

		// Parse Flat Transaction JSON
		var t struct {
			ID             int     `json:"id"`
			TotalAmount    float64 `json:"total_amount"`
			DiscountAmount float64 `json:"discount_amount"`
		}

		if err := json.Unmarshal(body, &t); err != nil {
			fmt.Println("❌ Error parsing JSON response:", err)
			fmt.Println("Raw Body:", string(body))
			return
		}

		fmt.Printf("   🧾 Transaction ID: %d\n", t.ID)
		fmt.Printf("   💰 Final Amount: %.2f\n", t.TotalAmount)
		fmt.Printf("   🏷️  Discount Amount: %.2f\n", t.DiscountAmount)
		fmt.Println("   -----------------------------")
		fmt.Printf("   Original Total: %.2f\n", t.TotalAmount+t.DiscountAmount)

		if t.DiscountAmount > 0 {
			fmt.Println("   🎉 DISCOUNT APPLIED CORRECTLY!")
		} else {
			fmt.Println("   ⚠️ No discount applied. Check if discount is active/period valid.")
		}

	} else {
		fmt.Println("❌ Checkout Failed:", resp.StatusCode)
		fmt.Println("Error:", string(body))
	}
}

func updateProductCategory(token string, productID, categoryID int) {
	fmt.Printf("Updating Product %d Category to %d...\n", productID, categoryID)
	// Get current product first to keep other fields
	// Using hardcoded values for simplicity in test
	product := map[string]interface{}{
		"nama":        fmt.Sprintf("Product %d (Test)", productID),
		"harga":       18000,
		"stok":        100,
		"category_id": categoryID,
	}
	jsonBody, _ := json.Marshal(product)

	// Using PUT endpoint
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/produk/%d", baseURL, productID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error updating product:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Product Category Updated!")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("⚠️ Update Product Failed:", resp.StatusCode, string(body))
	}
}
