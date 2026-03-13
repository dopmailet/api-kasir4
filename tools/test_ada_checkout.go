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

	// Ambil customer
	reqs, _ := http.NewRequest("GET", baseURL+"/api/customers?limit=1", nil)
	reqs.Header.Set("Authorization", auth)
	resps, _ := http.DefaultClient.Do(reqs)
	bodys, _ := io.ReadAll(resps.Body)
	resps.Body.Close()

	var lst map[string]interface{}
	json.Unmarshal(bodys, &lst)

	var customerID float64
	if dt, ok := lst["data"].([]interface{}); ok && len(dt) > 0 {
		cData := dt[0].(map[string]interface{})
		customerID = cData["id"].(float64)
	}
	if customerID == 0 {
		fmt.Println("❌ Tidak ada customer untuk ditest.")
		return
	}

	fmt.Printf("✅ Mencoba checkout untuk Customer ID: %.0f\n", customerID)

	// Ambil produk sembarang
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
		if idVal, ok := prod["id"].(float64); ok {
			productID = idVal
		}
		if priceVal, ok := prod["harga_jual"].(float64); ok {
			productPrice = priceVal
		}
	}

	actualTotal := productPrice * 1
	checkoutReq := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_id": int(productID),
				"quantity":   1,
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

	fmt.Println("Status Checkout:", respco.Status)
	fmt.Println("Body Checkout:", string(bodyco))
}
