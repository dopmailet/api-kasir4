//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api-kasir4dopmailet-production.up.railway.app"

func login() string {
	reqBody := map[string]string{
		"username": "tokoku",
		"password": "123456",
	}
	b, _ := json.Marshal(reqBody)
	resp, _ := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(b))
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var loginResp map[string]interface{}
	json.Unmarshal(body, &loginResp)

	token := ""
	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		token, _ = data["token"].(string)
	}
	return "Bearer " + token
}

func main() {
	auth := login()
	if auth == "Bearer " {
		fmt.Println("❌ Gagal login tokoku. Periksa username/password.")
		return
	}
	fmt.Println("✅ Login berhasil.")

	// Step 1: Create a temporary customer
	fmt.Println("\n--- 1. Bikin Customer Tester ---")
	custReq := map[string]interface{}{
		"name":  "Loyalty Tester",
		"phone": "081234567899",
	}
	bCust, _ := json.Marshal(custReq)
	reqc, _ := http.NewRequest("POST", baseURL+"/api/customers", bytes.NewBuffer(bCust))
	reqc.Header.Set("Authorization", auth)
	reqc.Header.Set("Content-Type", "application/json")
	respc, _ := http.DefaultClient.Do(reqc)
	bodyc, _ := io.ReadAll(respc.Body)
	respc.Body.Close()

	var cResp map[string]interface{}
	json.Unmarshal(bodyc, &cResp)

	// Try parsing customer ID
	var customerID float64
	if data, ok := cResp["data"].(map[string]interface{}); ok {
		customerID = data["id"].(float64)
	} else if idObj, ok := cResp["id"].(float64); ok {
		customerID = idObj
	} else if msg, ok := cResp["message"].(string); ok && msg == "Nomor telepon sudah digunakan" {
		fmt.Println("⚠️ Customer tester sudah ada di database, mencari manual...")
		// Ambil list customer
		reqs, _ := http.NewRequest("GET", baseURL+"/api/customers", nil)
		reqs.Header.Set("Authorization", auth)
		resps, _ := http.DefaultClient.Do(reqs)
		bodys, _ := io.ReadAll(resps.Body)
		resps.Body.Close()
		var lst map[string]interface{}
		json.Unmarshal(bodys, &lst)
		if dt, ok := lst["data"].([]interface{}); ok {
			for _, v := range dt {
				cData := v.(map[string]interface{})
				if cData["phone"].(string) == "081234567899" {
					customerID = cData["id"].(float64)
					break
				}
			}
		}
	}

	if customerID == 0 {
		fmt.Println("❌ Gagal mendapatkan ID customer tester.")
		fmt.Println(string(bodyc))
		return
	}
	fmt.Printf("✅ Customer ID: %.0f\n", customerID)

	// Step 2: Ambil produk sembarang
	fmt.Println("\n--- 2. Ambil Produk ---")
	reqp, _ := http.NewRequest("GET", baseURL+"/api/produk?limit=1", nil)
	reqp.Header.Set("Authorization", auth)
	respp, _ := http.DefaultClient.Do(reqp)
	bodyp, _ := io.ReadAll(respp.Body)
	respp.Body.Close()

	var pResp map[string]interface{}
	json.Unmarshal(bodyp, &pResp)

	var productID float64
	var productPrice float64
	if dt, ok := pResp["data"].([]interface{}); ok && len(dt) > 0 {
		prod := dt[0].(map[string]interface{})
		productID = prod["id"].(float64)
		productPrice = prod["harga_jual"].(float64)
	}

	if productID == 0 {
		fmt.Println("❌ Gagal menemukan produk untuk ditransaksikan.")
		return
	}
	fmt.Printf("✅ Product ID: %.0f (Harga: %.0f)\n", productID, productPrice)

	// Hitung quantity agar (qty * price) mendekati Rp 55.000,- (biar dapat 5 poin)
	targetTotal := 55000.0
	qty := int(targetTotal / productPrice)
	if qty < 1 {
		qty = 1
	}
	actualTotal := float64(qty) * productPrice

	// Step 3: Lakukan Checkout
	fmt.Printf("\n--- 3. Melakukan Checkout Rp %.0f ---\n", actualTotal)
	checkoutReq := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_id": int(productID),
				"quantity":   qty,
			},
		},
		"payment_amount": actualTotal,
		"customer_id":    int(customerID),
	}
	bCO, _ := json.Marshal(checkoutReq)
	reqco, _ := http.NewRequest("POST", baseURL+"/api/checkout", bytes.NewBuffer(bCO))
	reqco.Header.Set("Authorization", auth)
	reqco.Header.Set("Content-Type", "application/json")
	respco, _ := http.DefaultClient.Do(reqco)
	bodyco, _ := io.ReadAll(respco.Body)
	respco.Body.Close()

	if respco.StatusCode >= 300 {
		fmt.Println("❌ Checkout Gagal:")
		fmt.Println(string(bodyco))
		return
	}
	fmt.Println("✅ Checkout Sukses.")

	expectedPoints := int(actualTotal / 10000)
	fmt.Printf("   >> Harusnya dapat point: %d\n", expectedPoints)

	time.Sleep(2 * time.Second)

	// Step 4: Cek Profil Customer (Loyalty Points)
	fmt.Printf("\n--- 4. Cek Profil Customer (Loyalty Points) ---\n")
	reqGetC, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/customers/%.0f", baseURL, customerID), nil)
	reqGetC.Header.Set("Authorization", auth)
	respGetC, _ := http.DefaultClient.Do(reqGetC)
	bodyGetC, _ := io.ReadAll(respGetC.Body)
	respGetC.Body.Close()

	var profile map[string]interface{}
	json.Unmarshal(bodyGetC, &profile)

	actualPoints := -1.0
	if data, ok := profile["data"].(map[string]interface{}); ok {
		actualPoints = data["loyalty_points"].(float64)
	}

	if int(actualPoints) >= expectedPoints {
		fmt.Printf("✅ Poin customer bertambah dengan normal: %.0f point(s) total.\n", actualPoints)
	} else {
		fmt.Printf("❌ Poin customer FAIL. Diharapkan >= %d, tapi dapetnya %.0f\n", expectedPoints, actualPoints)
	}

	// Step 5: Cek Endpoint Riwayat Loyalty
	fmt.Printf("\n--- 5. Cek Riwayat Loyalty Transactions ---\n")
	reqH, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/customers/%.0f/loyalty-transactions", baseURL, customerID), nil)
	reqH.Header.Set("Authorization", auth)
	respH, _ := http.DefaultClient.Do(reqH)
	bodyH, _ := io.ReadAll(respH.Body)
	respH.Body.Close()

	fmt.Println(string(bodyH))

	// Validasi text simple
	if bytes.Contains(bodyH, []byte("earn")) {
		fmt.Println("✅ TEST PASS: Endpoint riwayat loyalty berhasil menampilkan 'earn'!")
	} else {
		fmt.Println("❌ TEST FAIL: Riwayat loyalty kosong / tidak ada 'earn'.")
	}
}
