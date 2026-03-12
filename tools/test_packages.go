//go:build ignore

// test_packages.go — Acceptance test untuk endpoint package superadmin
// Jalankan dengan: go run tools/test_packages.go
// (server harus berjalan di localhost:8080)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const baseURL = "http://localhost:8080/api"

var token string

func main() {
	log.Println("============================")
	log.Println(" Package Acceptance Tests")
	log.Println("============================")

	// Login sebagai superadmin (asumsi ada superadmin account)
	doLogin()

	var results []string
	pass := func(name string) { results = append(results, "✅ "+name) }
	fail := func(name, detail string) { results = append(results, "❌ "+name+" | "+detail) }

	// -------------------------------------------------------
	// Test 1: Edit fitur (tambah/hapus beberapa baris) → 200
	// -------------------------------------------------------
	log.Println("\n--- Test 1: Edit fitur (tambah/hapus) ---")
	pkgID := createTestPackage()
	if pkgID == 0 {
		fail("Test 1", "gagal buat test package")
	} else {
		body := makeBody(pkgID, "TestPkg", []interface{}{"Fitur A", "Fitur B", "Fitur C"})
		code, res := doRequest("PUT", fmt.Sprintf("/superadmin/packages/%d", pkgID), body)
		if code == 200 {
			features := extractFeatures(res)
			if len(features) == 3 {
				pass("Test 1: Edit fitur 3 baris")
			} else {
				fail("Test 1", fmt.Sprintf("features count = %d, want 3", len(features)))
			}
		} else {
			fail("Test 1", fmt.Sprintf("HTTP %d: %s", code, res))
		}
	}

	// -------------------------------------------------------
	// Test 2: Fitur dikosongkan semua → features = []
	// -------------------------------------------------------
	log.Println("\n--- Test 2: Fitur dikosongkan semua ---")
	if pkgID > 0 {
		body := makeBody(pkgID, "TestPkg", []interface{}{})
		code, res := doRequest("PUT", fmt.Sprintf("/superadmin/packages/%d", pkgID), body)
		if code == 200 {
			features := extractFeatures(res)
			if len(features) == 0 {
				pass("Test 2: Fitur kosong → [] tanpa 500")
			} else {
				fail("Test 2", fmt.Sprintf("features tidak kosong: %v", features))
			}
		} else {
			fail("Test 2", fmt.Sprintf("HTTP %d: %s", code, res))
		}
	}

	// -------------------------------------------------------
	// Test 3: Toggle aktif/nonaktif → 200
	// -------------------------------------------------------
	log.Println("\n--- Test 3: Toggle is_active ---")
	if pkgID > 0 {
		body := makeBodyFull(map[string]interface{}{
			"name": "TestPkg", "max_kasir": 1, "max_products": 100, "price": 50000,
			"is_active": false, "features": []string{}, "period": "/bulan",
			"discount_percent": 0, "is_popular": false,
		})
		code, res := doRequest("PUT", fmt.Sprintf("/superadmin/packages/%d", pkgID), body)
		if code == 200 {
			pass("Test 3: Toggle nonaktif → 200")
		} else {
			fail("Test 3", fmt.Sprintf("HTTP %d: %s", code, res))
		}
	}

	// -------------------------------------------------------
	// Test 4: Toggle populer → paket lain otomatis false
	// -------------------------------------------------------
	log.Println("\n--- Test 4: Toggle is_popular → reset paket lain ---")
	pkg2ID := createTestPackage()
	if pkg2ID > 0 && pkgID > 0 {
		// Set pkgID sebagai popular
		code, res := doRequest("PATCH", fmt.Sprintf("/superadmin/packages/%d/popular", pkgID),
			map[string]interface{}{"is_popular": true})
		if code != 200 {
			fail("Test 4", fmt.Sprintf("Gagal set popular pkgID: HTTP %d: %s", code, res))
		} else {
			// Cek pkg2 tidak populer
			_, getRes := doRequest("GET", fmt.Sprintf("/superadmin/packages/%d", pkg2ID), nil)
			var parsed map[string]interface{}
			json.Unmarshal([]byte(getRes), &parsed)
			if data, ok := parsed["data"].(map[string]interface{}); ok {
				if pop, _ := data["is_popular"].(bool); !pop {
					pass("Test 4: Paket lain otomatis false setelah toggle populer")
				} else {
					fail("Test 4", "pkg2 masih true padahal pkgID yang dipopulerkan")
				}
			} else {
				fail("Test 4", "Gagal parse GET response pkg2")
			}
		}
	}

	// -------------------------------------------------------
	// Test 5: description dan discount_label kosong → NULL
	// -------------------------------------------------------
	log.Println("\n--- Test 5: Nullable fields → NULL ---")
	if pkgID > 0 {
		body := makeBodyFull(map[string]interface{}{
			"name": "TestPkg", "max_kasir": 1, "max_products": 100, "price": 50000,
			"is_active": true, "features": []string{}, "period": "/bulan",
			"discount_percent": 0, "is_popular": false,
			"description": "", "discount_label": "",
		})
		code, res := doRequest("PUT", fmt.Sprintf("/superadmin/packages/%d", pkgID), body)
		if code == 200 {
			var parsed map[string]interface{}
			json.Unmarshal([]byte(res), &parsed)
			if data, ok := parsed["data"].(map[string]interface{}); ok {
				descNull := data["description"] == nil
				labelNull := data["discount_label"] == nil
				if descNull && labelNull {
					pass("Test 5: description dan discount_label → null")
				} else {
					fail("Test 5", fmt.Sprintf("description=%v discount_label=%v (want null null)", data["description"], data["discount_label"]))
				}
			}
		} else {
			fail("Test 5", fmt.Sprintf("HTTP %d: %s", code, res))
		}
	}

	// -------------------------------------------------------
	// Test 6: Update harga/periode/discount_percent → 200
	// -------------------------------------------------------
	log.Println("\n--- Test 6: Update harga/periode/discount ---")
	if pkgID > 0 {
		body := makeBodyFull(map[string]interface{}{
			"name": "TestPkg-Updated", "max_kasir": 5, "max_products": 500, "price": 299000,
			"is_active": true, "features": []string{"Fitur X"}, "period": "/tahun",
			"discount_percent": 10, "discount_label": "Hemat!", "is_popular": false,
		})
		code, _ := doRequest("PUT", fmt.Sprintf("/superadmin/packages/%d", pkgID), body)
		if code == 200 {
			pass("Test 6: Update harga/periode/discount → 200")
		} else {
			fail("Test 6", fmt.Sprintf("HTTP %d", code))
		}
	}

	// -------------------------------------------------------
	// Test 7: POST create paket fitur kosong + nullable kosong
	// -------------------------------------------------------
	log.Println("\n--- Test 7: POST create paket minimal ---")
	body7 := makeBodyFull(map[string]interface{}{
		"name": "PaketTest7", "max_kasir": 1, "max_products": 50, "price": 0,
		"is_active": false, "features": []interface{}{}, "period": "/bulan",
		"discount_percent": 0, "is_popular": false,
	})
	code7, res7 := doRequest("POST", "/superadmin/packages", body7)
	if code7 == 201 {
		pass("Test 7: POST create paket minimal → 201")
		// Cleanup
		var parsed map[string]interface{}
		json.Unmarshal([]byte(res7), &parsed)
		if data, ok := parsed["data"].(map[string]interface{}); ok {
			if newID, ok := data["id"].(float64); ok {
				doRequest("DELETE", fmt.Sprintf("/superadmin/packages/%d", int(newID)), nil)
			}
		}
	} else {
		fail("Test 7", fmt.Sprintf("HTTP %d: %s", code7, res7))
	}

	// -------------------------------------------------------
	// Test 8: Format features sebagai []object → normalisasi
	// -------------------------------------------------------
	log.Println("\n--- Test 8: features format []object ---")
	if pkgID > 0 {
		body8 := makeBodyFull(map[string]interface{}{
			"name": "TestPkg", "max_kasir": 1, "max_products": 100, "price": 50000,
			"is_active": true, "period": "/bulan", "discount_percent": 0, "is_popular": false,
			"features": []map[string]string{{"label": "Fitur Object A"}, {"label": "Fitur Object B"}},
		})
		code, res := doRequest("PUT", fmt.Sprintf("/superadmin/packages/%d", pkgID), body8)
		if code == 200 {
			features := extractFeatures(res)
			if len(features) == 2 && strings.Contains(features[0], "Fitur Object") {
				pass("Test 8: features []object → normalisasi ke []string")
			} else {
				fail("Test 8", fmt.Sprintf("features=%v", features))
			}
		} else {
			fail("Test 8", fmt.Sprintf("HTTP %d: %s", code, res))
		}
	}

	// Cleanup: hapus test packages
	if pkgID > 0 {
		doRequest("DELETE", fmt.Sprintf("/superadmin/packages/%d", pkgID), nil)
	}
	if pkg2ID > 0 {
		doRequest("DELETE", fmt.Sprintf("/superadmin/packages/%d", pkg2ID), nil)
	}

	// Print results
	log.Println("\n============================")
	log.Println("  HASIL ACCEPTANCE TESTS")
	log.Println("============================")
	for _, r := range results {
		log.Println(r)
	}

	passCount := 0
	for _, r := range results {
		if strings.HasPrefix(r, "✅") {
			passCount++
		}
	}
	log.Printf("\nTotal: %d/%d passed\n", passCount, len(results))
}

func doLogin() {
	payload := map[string]string{"username": "superadmin", "password": "SuperAdmin123!"}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal("Gagal login:", err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if data, ok := res["data"].(map[string]interface{}); ok {
		if t, ok := data["token"].(string); ok {
			token = t
			log.Println("✅ Login superadmin berhasil")
			return
		}
	}
	log.Fatal("❌ Gagal dapat token superadmin:", res)
}

func doRequest(method, path string, body interface{}) (int, string) {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	}
	req, _ := http.NewRequest(method, baseURL+path, reqBody)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)
	log.Printf("  → %s %s | HTTP %d | %s", method, path, resp.StatusCode, string(rawBody))
	return resp.StatusCode, string(rawBody)
}

func createTestPackage() int {
	body := makeBodyFull(map[string]interface{}{
		"name": "TestPkgTemp", "max_kasir": 1, "max_products": 100, "price": 10000,
		"is_active": true, "features": []string{"Fitur Init"}, "period": "/bulan",
		"discount_percent": 0, "is_popular": false,
	})
	code, res := doRequest("POST", "/superadmin/packages", body)
	if code != 201 {
		log.Printf("createTestPackage gagal: HTTP %d %s", code, res)
		return 0
	}
	var parsed map[string]interface{}
	json.Unmarshal([]byte(res), &parsed)
	if data, ok := parsed["data"].(map[string]interface{}); ok {
		if id, ok := data["id"].(float64); ok {
			return int(id)
		}
	}
	return 0
}

func makeBody(id int, name string, features []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name": name, "max_kasir": 1, "max_products": 100, "price": 50000,
		"is_active": true, "features": features, "period": "/bulan",
		"discount_percent": 0, "is_popular": false,
	}
}

func makeBodyFull(m map[string]interface{}) map[string]interface{} {
	return m
}

func extractFeatures(res string) []string {
	var parsed map[string]interface{}
	json.Unmarshal([]byte(res), &parsed)
	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		return nil
	}
	raw, ok := data["features"].([]interface{})
	if !ok {
		return []string{}
	}
	result := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
