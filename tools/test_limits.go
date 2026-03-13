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

func main() {
	client := &http.Client{Timeout: 15 * time.Second}

	// =========================================================
	// STEP 1: Login
	// =========================================================
	fmt.Println("\n============================")
	fmt.Println("STEP 1: LOGIN")
	fmt.Println("============================")

	loginBody := `{"username":"tokoku","password":"123456"}`
	resp, err := client.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBufferString(loginBody))
	if err != nil {
		fmt.Printf("❌ Login gagal: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var loginResp map[string]interface{}
	json.Unmarshal(body, &loginResp)

	// Response format: {"data": {"token": "...", "user": {...}}}
	token := ""
	var user map[string]interface{}

	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		token, _ = data["token"].(string)
		user, _ = data["user"].(map[string]interface{})
	}
	// Fallback: token langsung di root (format lama)
	if token == "" {
		token, _ = loginResp["token"].(string)
		user, _ = loginResp["user"].(map[string]interface{})
	}

	if token == "" {
		fmt.Printf("❌ Token tidak ditemukan. Response: %s\n", string(body))
		return
	}
	fmt.Printf("✅ Login berhasil!\n")
	fmt.Printf("   Username : %v\n", user["username"])
	fmt.Printf("   Role     : %v\n", user["role"])
	fmt.Printf("   StoreID  : %v\n", user["store_id"])
	fmt.Printf("   Token    : %s...\n", token[:30])

	authHeader := "Bearer " + token
	today := time.Now().In(mustLoadLoc("Asia/Makassar")).Format("2006-01-02")

	// =========================================================
	// STEP 2: GET /api/store/limits
	// =========================================================
	fmt.Println("\n============================")
	fmt.Println("STEP 2: GET /api/store/limits")
	fmt.Println("============================")

	limitsURL := baseURL + "/api/store/limits?timezone=Asia/Makassar"
	limitsData := doGet(client, limitsURL, authHeader)

	todaySalesBefore := 0
	if data, ok := limitsData["data"].(map[string]interface{}); ok {
		prettyPrint(data)
		if ts, ok := data["today_sales"].(float64); ok {
			todaySalesBefore = int(ts)
		}
	} else {
		fmt.Println("Raw response:", limitsData)
	}
	fmt.Printf("\n📊 today_sales SEBELUM checkout: %d\n", todaySalesBefore)

	// =========================================================
	// STEP 3: GET /api/dashboard/summary (pembanding)
	// =========================================================
	fmt.Println("\n============================")
	fmt.Println("STEP 3: GET /api/dashboard/summary (pembanding)")
	fmt.Println("============================")

	summaryURL := fmt.Sprintf("%s/api/dashboard/summary?start_date=%s&end_date=%s&timezone=Asia/Makassar", baseURL, today, today)
	summaryData := doGet(client, summaryURL, authHeader)

	totalTransaksi := 0
	if data, ok := summaryData["data"].(map[string]interface{}); ok {
		prettyPrint(data)
		if tt, ok := data["total_transaksi"].(float64); ok {
			totalTransaksi = int(tt)
		}
	} else {
		prettyPrint(summaryData)
	}
	fmt.Printf("\n📊 total_transaksi dari dashboard/summary: %d\n", totalTransaksi)

	// Bandingkan
	fmt.Println("\n============================")
	fmt.Println("PERBANDINGAN")
	fmt.Println("============================")
	if todaySalesBefore == totalTransaksi {
		fmt.Printf("✅ MATCH! today_sales (%d) == total_transaksi (%d)\n", todaySalesBefore, totalTransaksi)
	} else {
		fmt.Printf("❌ MISMATCH! today_sales (%d) != total_transaksi (%d)\n", todaySalesBefore, totalTransaksi)
	}

	// =========================================================
	// STEP 4: GET /api/produk (ambil produk pertama)
	// =========================================================
	fmt.Println("\n============================")
	fmt.Println("STEP 4: GET /api/produk (ambil produk pertama)")
	fmt.Println("============================")

	produkData := doGet(client, baseURL+"/api/produk?limit=1", authHeader)
	productID := 0
	productPrice := 0.0
	productName := ""

	// Coba ambil dari berbagai format response
	if dataArr, ok := produkData["data"].([]interface{}); ok && len(dataArr) > 0 {
		if p, ok := dataArr[0].(map[string]interface{}); ok {
			if id, ok := p["id"].(float64); ok {
				productID = int(id)
			}
			if hp, ok := p["harga_jual"].(float64); ok {
				productPrice = hp
			}
			if nm, ok := p["nama"].(string); ok {
				productName = nm
			}
		}
	}

	if productID == 0 {
		fmt.Println("❌ Tidak ada produk tersedia, skip checkout test")
		return
	}
	fmt.Printf("✅ Produk ditemukan: ID=%d, Nama=%s, Harga=%.0f\n", productID, productName, productPrice)

	// =========================================================
	// STEP 5: POST /api/checkout
	// =========================================================
	fmt.Println("\n============================")
	fmt.Println("STEP 5: POST /api/checkout")
	fmt.Println("============================")

	checkoutBody := fmt.Sprintf(`{
		"items": [{"product_id": %d, "quantity": 1, "price": %.0f}],
		"payment_amount": %.0f,
		"timezone": "Asia/Makassar"
	}`, productID, productPrice, productPrice)

	req, _ := http.NewRequest("POST", baseURL+"/api/checkout?timezone=Asia/Makassar", bytes.NewBufferString(checkoutBody))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	checkoutResp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Checkout gagal: %v\n", err)
		return
	}
	defer checkoutResp.Body.Close()
	checkoutBody2, _ := io.ReadAll(checkoutResp.Body)
	var checkoutResult map[string]interface{}
	json.Unmarshal(checkoutBody2, &checkoutResult)
	fmt.Printf("HTTP Status: %d\n", checkoutResp.StatusCode)
	if checkoutResp.StatusCode == 201 {
		fmt.Printf("✅ Checkout berhasil! Transaction ID: %v\n", checkoutResult["id"])
	} else {
		fmt.Printf("❌ Checkout gagal: %s\n", string(checkoutBody2))
		return
	}

	// =========================================================
	// STEP 6: GET /api/store/limits lagi (cek increment)
	// =========================================================
	fmt.Println("\n============================")
	fmt.Println("STEP 6: GET /api/store/limits (setelah checkout)")
	fmt.Println("============================")

	time.Sleep(1 * time.Second) // beri waktu server update
	limitsData2 := doGet(client, limitsURL, authHeader)

	todaySalesAfter := 0
	if data, ok := limitsData2["data"].(map[string]interface{}); ok {
		prettyPrint(data)
		if ts, ok := data["today_sales"].(float64); ok {
			todaySalesAfter = int(ts)
		}
	}

	fmt.Println("\n============================")
	fmt.Println("HASIL FINAL")
	fmt.Println("============================")
	fmt.Printf("today_sales SEBELUM : %d\n", todaySalesBefore)
	fmt.Printf("today_sales SESUDAH : %d\n", todaySalesAfter)
	if todaySalesAfter == todaySalesBefore+1 {
		fmt.Println("✅ TEST PASS! Transaksi tercatat dengan benar (+1)")
	} else {
		fmt.Printf("❌ TEST FAIL! Expected %d, got %d\n", todaySalesBefore+1, todaySalesAfter)
	}
}

func doGet(client *http.Client, url, auth string) map[string]interface{} {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error GET %s: %v\n", url, err)
		return nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result
}

func prettyPrint(data map[string]interface{}) {
	for k, v := range data {
		fmt.Printf("   %-25s: %v\n", k, v)
	}
}

func mustLoadLoc(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}
