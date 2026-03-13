package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	fmt.Println("Token:", auth[:20]+"...")

	req, _ := http.NewRequest("GET", baseURL+"/api/cash-flow/funds?limit=10&page=1", nil)
	req.Header.Set("Authorization", auth)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("RESPONSE /api/cash-flow/funds:")
	fmt.Println(string(body))
}
