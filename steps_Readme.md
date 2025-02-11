# Receipt Processor Web Service

This is a web service for processing receipts and calculating reward points based on predefined rules.

## Steps to Run the Web Service

### 1. Clone the Repository
```sh
git clone https://github.com/YOUR_REPOSITORY/receipt-processor-challenge.git
cd receipt-processor-challenge
```

### 2. Build the Docker Image
```sh
docker build -t receipt_processor .
```

### 3. Run the Docker Container
```sh
docker run -dp 8080:8080 receipt_processor
```

### 4. Process a Receipt Using cURL (this command is for CMD)
Submit a receipt to be processed and receive an ID in response:
```sh
curl -X POST "http://localhost:8080/receipts/process" ^
     -H "Content-Type: application/json" ^
     -d "{ \"retailer\": \"Target\", \"purchaseDate\": \"2022-01-01\", \"purchaseTime\": \"13:01\", \"total\": 35.35, \"items\": [ { \"shortDescription\": \"Mountain Dew 12PK\", \"price\": 6.49 }, { \"shortDescription\": \"Emils Cheese Pizza\", \"price\": 12.25 }, { \"shortDescription\": \"Knorr Creamy Chicken\", \"price\": 1.26 }, { \"shortDescription\": \"Doritos Nacho Cheese\", \"price\": 3.35 }, { \"shortDescription\": \"Klarbrunn 12-PK 12 FL OZ\", \"price\": 12.00 } ] }"

```

Example response:
```json
{"id":"your-id-here"}
```

### 5. Retrieve Points for a Processed Receipt
Use the ID from the previous step to get the points assigned to the receipt:
```sh
curl -X GET http://localhost:8080/receipts/your-id-here/points
```

Example response:
```json
{"points": 28}
```

## Notes
- Ensure that port `8080` is free before running the service.
- Data is stored in memory and will be lost when the service restarts.
- You can use Postman or any API testing tool instead of cURL.

This web service is designed to meet the Receipt Processor Challenge API requirements.

