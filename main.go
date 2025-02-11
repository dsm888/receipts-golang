package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Receipt struct to store receipt details
// Includes information such as retailer name, purchase date/time, total amount, items, and calculated points
type Receipt struct {
	ID           string  `json:"id"`
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        float64 `json:"total"`
	Items        []Item  `json:"items"`
	Points       int     `json:"points"`
}

// Item struct represents an item on the receipt
// Each item has a short description and a price
type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

// Global map to store receipts using their ID as the key
var receipts map[string]Receipt

// init function initializes the global receipts map
func init() {
	receipts = make(map[string]Receipt)
}

func main() {
	// Create a new router using Gorilla Mux
	router := mux.NewRouter()

	// Define routes for processing receipts and retrieving points
	router.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	// Start the HTTP server on port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}

// ProcessReceipt handles the POST request to process a receipt
// It decodes the JSON request, validates the data, calculates points, and stores the receipt
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate date and time format
	if len(receipt.PurchaseDate) != 10 || len(receipt.PurchaseTime) != 5 {
		http.Error(w, "Invalid date or time format", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the receipt
	receipt.ID = uuid.New().String()
	receipt.Points = calculatePoints(receipt)
	receipts[receipt.ID] = receipt

	// Respond with the generated receipt ID
	response := map[string]string{"id": receipt.ID}
	json.NewEncoder(w).Encode(response)
}

// GetPoints handles the GET request to retrieve points for a given receipt ID
func GetPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	receipt, found := receipts[id]
	if !found {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Respond with the points associated with the receipt
	response := map[string]int{"points": receipt.Points}
	json.NewEncoder(w).Encode(response)
}

// calculatePoints calculates the reward points based on receipt details
func calculatePoints(receipt Receipt) int {
	points := 0

	retailerName := strings.ReplaceAll(receipt.Retailer, "&", "") // Remove special chars
	retailerName = strings.ReplaceAll(retailerName, " ", "")      // Remove spaces
	points += len(retailerName)                                   // Count only alphanumeric chars

	if isRoundDollar(receipt.Total) {
		points += 50
	}

	if isMultipleOfQuarter(receipt.Total) {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	purchaseDay, _ := strconv.Atoi(receipt.PurchaseDate[len(receipt.PurchaseDate)-2:])
	if purchaseDay%2 != 0 {
		points += 6
	}

	purchaseHour, _ := strconv.Atoi(strings.Split(receipt.PurchaseTime, ":")[0])
	if purchaseHour >= 14 && purchaseHour < 16 {
		points += 10
	}

	return points
}

// isRoundDollar checks if the total amount is a whole dollar
func isRoundDollar(amount float64) bool {
	return amount == float64(int(amount))
}

// isMultipleOfQuarter checks if the total is a multiple of $0.25
func isMultipleOfQuarter(amount float64) bool {
	return math.Mod((amount*100), 25.0) == 0.0
}
