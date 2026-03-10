package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	baseURL := "http://localhost:8080/api"

	// 1. Login as Admin
	loginData := map[string]string{"username": "admin", "password": "admin123"}
	loginBytes, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBytes))
	if err != nil {
		fmt.Println("Error login:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var loginResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &loginResp); err != nil {
		fmt.Println("Error unmarshal login:", err)
		os.Exit(1)
	}
	token := loginResp.Data.Token
	if token == "" {
		fmt.Println("No token, login failed:", string(body))
		os.Exit(1)
	}
	fmt.Println("Logged in successfully.")

	client := &http.Client{}

	// Optional: create a purchase to test it
	// But let's just get the ledger directly first
	reqLedger, _ := http.NewRequest("GET", baseURL+"/cash-flow/ledger", nil)
	reqLedger.Header.Set("Authorization", "Bearer "+token)
	ledgResp, err := client.Do(reqLedger)
	if err != nil {
		fmt.Println("Error getting ledger:", err)
		os.Exit(1)
	}
	defer ledgResp.Body.Close()
	ledgBody, _ := ioutil.ReadAll(ledgResp.Body)

	fmt.Println("\n--- LEDGER RESPONSE ---")
	fmt.Println(string(ledgBody))

	// Get summary to compare
	reqSummary, _ := http.NewRequest("GET", baseURL+"/cash-flow/summary", nil)
	reqSummary.Header.Set("Authorization", "Bearer "+token)
	sumResp, _ := client.Do(reqSummary)
	if err != nil {
		fmt.Println("Error getting summary:", err)
		os.Exit(1)
	}
	defer sumResp.Body.Close()
	sumBody, _ := ioutil.ReadAll(sumResp.Body)

	fmt.Println("\n--- SUMMARY RESPONSE ---")
	fmt.Println(string(sumBody))
}
