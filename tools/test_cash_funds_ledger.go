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

func login(username, password string) (string, error) {
	reqBody := map[string]string{
		"username": username,
		"password": password,
	}
	b, _ := json.Marshal(reqBody)

	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var loginResp map[string]interface{}
	json.Unmarshal(body, &loginResp)

	// Format nested { data: { token: ... } }
	token := ""
	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		token, _ = data["token"].(string)
	}
	if token == "" {
		token, _ = loginResp["token"].(string)
	}

	if token == "" {
		return "", fmt.Errorf("token not found in response: %s", string(body))
	}
	return "Bearer " + token, nil
}

func getInitialBalance(auth string) (float64, error) {
	req, _ := http.NewRequest("GET", baseURL+"/api/cash-flow/initial-balance", nil)
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var respData map[string]interface{}
	json.Unmarshal(body, &respData)

	if amount, ok := respData["amount"].(float64); ok {
		return amount, nil
	}
	return 0, fmt.Errorf("failed parsing amount: %s", string(body))
}

func getSummary(auth string) (float64, error) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/cash-flow/summary?start_date=%s&end_date=%s", baseURL, today, today), nil)
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var respData map[string]interface{}
	json.Unmarshal(body, &respData)

	if data, ok := respData["data"].(map[string]interface{}); ok {
		if sk, ok := data["saldo_kas"].(float64); ok {
			return sk, nil
		}
	}
	return 0, fmt.Errorf("failed parsing saldo_kas: %s", string(body))
}

func getLedger(auth string) ([]interface{}, error) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/cash-flow/ledger?start_date=%s&end_date=%s&limit=20", baseURL, today, today), nil)
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var respData map[string]interface{}
	json.Unmarshal(body, &respData)

	if data, ok := respData["data"].([]interface{}); ok {
		return data, nil
	}
	return nil, fmt.Errorf("failed parsing ledger data: %s", string(body))
}

func main() {
	fmt.Println("======================================")
	fmt.Println("TEST 1: Skenario Ledger & Initial Balance (STORE 6)")
	fmt.Println("======================================")

	authTokoku, err := login("tokoku", "123456")
	if err != nil {
		fmt.Println("❌ Login tokoku gagal:", err)
		return
	}
	fmt.Println("✅ Login tokoku berhasil!")

	initBalBefore, _ := getInitialBalance(authTokoku)
	summaryBefore, _ := getSummary(authTokoku)
	fmt.Printf("📊 SEBELUM -> Initial Balance: %.0f | Saldo Kas: %.0f\n", initBalBefore, summaryBefore)

	// Test Fund In
	fmt.Println("\n--- POST /api/cash-flow/funds (IN: 1000000) ---")
	reqBodyIn := map[string]interface{}{
		"type":        "in",
		"amount":      1000000,
		"date":        time.Now().Format("2006-01-02"),
		"description": "Modal Awal Tokoku",
	}
	bIn, _ := json.Marshal(reqBodyIn)
	reqIn, _ := http.NewRequest("POST", baseURL+"/api/cash-flow/funds", bytes.NewBuffer(bIn))
	reqIn.Header.Set("Content-Type", "application/json")
	reqIn.Header.Set("Authorization", authTokoku)
	http.DefaultClient.Do(reqIn)
	fmt.Println("✅ Modal masuk (IN) tercatat.")

	// Test Fund Out
	fmt.Println("\n--- POST /api/cash-flow/funds (OUT: 200000) ---")
	reqBodyOut := map[string]interface{}{
		"type":        "out",
		"amount":      200000,
		"date":        time.Now().Format("2006-01-02"),
		"description": "Tarik Tunai Modal",
	}
	bOut, _ := json.Marshal(reqBodyOut)
	reqOut, _ := http.NewRequest("POST", baseURL+"/api/cash-flow/funds", bytes.NewBuffer(bOut))
	reqOut.Header.Set("Content-Type", "application/json")
	reqOut.Header.Set("Authorization", authTokoku)
	http.DefaultClient.Do(reqOut)
	fmt.Println("✅ Modal keluar (OUT) tercatat.")

	time.Sleep(2 * time.Second)

	initBalAfter, _ := getInitialBalance(authTokoku)
	summaryAfter, _ := getSummary(authTokoku)
	fmt.Printf("📊 SESUDAH -> Initial Balance: %.0f | Saldo Kas: %.0f\n", initBalAfter, summaryAfter)

	expectedDiff := float64(800000) // 1000000 - 200000
	if initBalAfter-initBalBefore == expectedDiff {
		fmt.Println("✅ TEST PASS: Initial Balance bertambah tepat 800.000")
	} else {
		fmt.Println("❌ TEST FAIL: Initial Balance tidak sesuai harapan")
	}
	if summaryAfter-summaryBefore == expectedDiff {
		fmt.Println("✅ TEST PASS: Saldo Kas (Summary) bertambah tepat 800.000")
	} else {
		fmt.Println("❌ TEST FAIL: Saldo Kas (Summary) tidak sesuai harapan")
	}

	fmt.Println("\n--- Cek Ledger ---")
	ledgerData, err := getLedger(authTokoku)
	if err != nil {
		fmt.Println("Kesalahan:", err)
	} else {
		fundInFound := false
		fundOutFound := false
		fmt.Printf("Ditemukan %d entri di halaman 1.\n", len(ledgerData))
		for _, entry := range ledgerData {
			e := entry.(map[string]interface{})
			cat := e["category"].(string)
			amt := e["amount"].(float64)
			rb := e["running_balance"].(float64)
			fmt.Printf("   -> %s (%.0f) | RunBal: %.0f | %s\n", cat, amt, rb, e["description"])
			if cat == "fund_in" {
				fundInFound = true
			}
			if cat == "fund_out" {
				fundOutFound = true
			}
		}
		if fundInFound && fundOutFound {
			fmt.Println("✅ TEST PASS: Kategori 'fund_in' dan 'fund_out' muncul di Ledger!")
		} else {
			fmt.Println("❌ TEST FAIL: Kategori fund tidak lengkap di Ledger!")
		}
	}

	fmt.Println("\n======================================")
	fmt.Println("TEST 2: ISOLASI MULTI-TENANT (STORE 5 - akhzayn)")
	fmt.Println("======================================")

	authAkhzayn, err := login("akhzayn", "123456")
	if err != nil {
		fmt.Println("❌ Login akhzayn gagal:", err)
		return
	}
	fmt.Println("✅ Login akhzayn berhasil!")

	initBalAkhzayn, _ := getInitialBalance(authAkhzayn)
	fmt.Printf("📊 Initial Balance Akhzayn: %.0f\n", initBalAkhzayn)

	ledgerAkhzayn, _ := getLedger(authAkhzayn)
	bocor := false
	for _, entry := range ledgerAkhzayn {
		e := entry.(map[string]interface{})
		desc := e["description"].(string)
		if desc == "Modal Awal Tokoku" || desc == "Tarik Tunai Modal" {
			bocor = true
		}
	}

	if bocor {
		fmt.Println("❌ TEST FAIL: Data tokoku BOCOR masuk ke ledger akhzayn!")
	} else {
		fmt.Println("✅ TEST PASS: Isolasi berhasil. Data tokoku tidak ada di ledger akhzayn.")
	}

	fmt.Println("\nSelesai.")
}
