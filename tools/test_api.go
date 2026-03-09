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
			User struct {
				ID int `json:"id"`
			} `json:"user"`
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(body, &loginResp)
	token := loginResp.Data.Token
	userID := loginResp.Data.User.ID
	fmt.Println("Logged in. Admin ID:", userID)

	// 2. Test Get Users
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	usersResp, _ := client.Do(req)
	defer usersResp.Body.Close()
	usersBody, _ := ioutil.ReadAll(usersResp.Body)
	fmt.Println("GET /users:", string(usersBody)[:min(200, len(string(usersBody)))], "...")

	// 3. Test Transactions with user_id mapping
	reqTx, _ := http.NewRequest("GET", fmt.Sprintf("%s/transactions?user_id=%d", baseURL, userID), nil)
	reqTx.Header.Set("Authorization", "Bearer "+token)
	txResp, _ := client.Do(reqTx)
	defer txResp.Body.Close()
	txBody, _ := ioutil.ReadAll(txResp.Body)
	fmt.Println("GET /transactions?user_id:", string(txBody)[:min(200, len(string(txBody)))], "...")

	// 4. Test Report Filter with user_id mapping
	reqRep, _ := http.NewRequest("GET", fmt.Sprintf("%s/report/hari-ini?user_id=%d", baseURL, userID), nil)
	reqRep.Header.Set("Authorization", "Bearer "+token)
	repResp, _ := client.Do(reqRep)
	defer repResp.Body.Close()
	repBody, _ := ioutil.ReadAll(repResp.Body)
	fmt.Println("GET /report/hari-ini?user_id:", string(repBody)[:min(200, len(string(repBody)))], "...")

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
