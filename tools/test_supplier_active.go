package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api-kasir4dopmailet-production.up.railway.app"

func loginADA() string {
	reqBody := map[string]string{
		"username": "ada",
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
	auth := loginADA()
	if auth == "Bearer " {
		fmt.Println("❌ Gagal login 'ada'. Periksa username/password.")
		return
	}
	fmt.Println("✅ Login 'ada' berhasil.")

	// 1. Buat supplier baru tanpa mengirim is_active
	fmt.Println("\n--- 1. Testing CREATE Supplier (Default is_active = true) ---")
	newSupplierReq := map[string]interface{}{
		"name":  "Supplier Test Auto-Active",
		"phone": "08999999999",
	}
	b, _ := json.Marshal(newSupplierReq)
	reqc, _ := http.NewRequest("POST", baseURL+"/api/suppliers", bytes.NewBuffer(b))
	reqc.Header.Set("Authorization", auth)
	reqc.Header.Set("Content-Type", "application/json")
	respc, _ := http.DefaultClient.Do(reqc)
	bodyc, _ := io.ReadAll(respc.Body)
	respc.Body.Close()

	var cResp map[string]interface{}
	json.Unmarshal(bodyc, &cResp)

	var supplierID float64
	var isActive bool
	if data, ok := cResp["data"].(map[string]interface{}); ok {
		supplierID = data["id"].(float64)
		isActive = data["is_active"].(bool)
		fmt.Printf("✅ Supplier berhasil dibuat. ID: %.0f, is_active: %v\n", supplierID, isActive)
	} else {
		fmt.Println("❌ Gagal membuat supplier:", string(bodyc))
		return
	}

	if !isActive {
		fmt.Println("❌ ERROR: is_active seharusnya TRUE secara default!")
	}

	// 2. Disable supplier (is_active = false)
	fmt.Println("\n--- 2. Testing UPDATE Supplier (Set is_active = false) ---")
	updateSupplierReq := map[string]interface{}{
		"is_active": false,
	}
	bu, _ := json.Marshal(updateSupplierReq)
	requ, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/suppliers/%.0f", baseURL, supplierID), bytes.NewBuffer(bu))
	requ.Header.Set("Authorization", auth)
	requ.Header.Set("Content-Type", "application/json")
	respu, _ := http.DefaultClient.Do(requ)
	bodyu, _ := io.ReadAll(respu.Body)
	respu.Body.Close()

	var uResp map[string]interface{}
	json.Unmarshal(bodyu, &uResp)

	if data, ok := uResp["data"].(map[string]interface{}); ok {
		isActive = data["is_active"].(bool)
		fmt.Printf("✅ Supplier di-update. is_active sekarang: %v\n", isActive)
	} else {
		fmt.Println("❌ Gagal update supplier:", string(bodyu))
		return
	}

	if isActive {
		fmt.Println("❌ ERROR: is_active seharusnya FALSE setelah update!")
	}

	// 3. GET All suppliers
	fmt.Println("\n--- 3. Testing GET API untuk visibilitas ---")
	reqg, _ := http.NewRequest("GET", baseURL+"/api/suppliers?search=Supplier Test Auto-Active", nil)
	reqg.Header.Set("Authorization", auth)
	respg, _ := http.DefaultClient.Do(reqg)
	bodyg, _ := io.ReadAll(respg.Body)
	respg.Body.Close()

	var gResp map[string]interface{}
	json.Unmarshal(bodyg, &gResp)

	if dt, ok := gResp["data"].([]interface{}); ok && len(dt) > 0 {
		for _, item := range dt {
			sup := item.(map[string]interface{})
			if sup["id"].(float64) == supplierID {
				fmt.Printf("✅ Supplier ID %.0f ditemukan di GET /api/suppliers. is_active: %v\n", supplierID, sup["is_active"])
			}
		}
	} else {
		fmt.Println("❌ Gagal fetch supplier list:", string(bodyg))
	}

	fmt.Println("\n✅ SEMUA TEST SELESAI. Silakan login ke Dashboard Frontend dan pastikan dropdown supplier di Pembelian Baru bekerja sesuai filter.")
}
