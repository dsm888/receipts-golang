package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestProcessReceipt(t *testing.T) {
	reqBody := `{
		"retailer": "M&M Corner Market",
		"purchaseDate": "2022-03-20",
		"purchaseTime": "14:33",
		"total": 23.24,
		"items": [
			{ "shortDescription": "Gatorade", "price": 2.25 },
			{ "shortDescription": "Mountain Dew 12PK", "price": 6.49 },
			{ "shortDescription": "Emils Cheese Pizza", "price": 12.25 },
			{ "shortDescription": "Pepsi", "price": 2.25 }
		]
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ProcessReceipt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if _, exists := resp["id"]; !exists {
		t.Errorf("Expected response to contain 'id'")
	}
}

func TestGetPoints(t *testing.T) {
	// Set up router
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	// First, create a receipt
	reqBody := `{
		"retailer": "M&M Corner Market",
		"purchaseDate": "2022-03-20",
		"purchaseTime": "14:33",
		"total": 23.24,
		"items": [
			{ "shortDescription": "Gatorade", "price": 2.25 },
			{ "shortDescription": "Mountain Dew 12PK", "price": 6.49 },
			{ "shortDescription": "Emils Cheese Pizza", "price": 12.25 },
			{ "shortDescription": "Pepsi", "price": 2.25 }
		]
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	id, exists := resp["id"]
	if !exists {
		t.Fatalf("Failed to obtain receipt ID")
	}

	// Now, retrieve the points
	pointsReq := httptest.NewRequest("GET", "/receipts/"+id+"/points", nil)
	pointsW := httptest.NewRecorder()
	router.ServeHTTP(pointsW, pointsReq)

	if pointsW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", pointsW.Code)
	}

	var pointsResp map[string]int
	json.Unmarshal(pointsW.Body.Bytes(), &pointsResp)
	if _, exists := pointsResp["points"]; !exists {
		t.Errorf("Expected response to contain 'points'")
	}
}
