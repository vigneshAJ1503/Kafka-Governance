# API Status Codes and Error Handling

## Overview
All API endpoints now have comprehensive status code handling with proper validation and error responses.

---

## HTTP Status Codes Used

| Code | Meaning | Usage |
|------|---------|-------|
| **201** | Created | Successful resource creation (POST) |
| **200** | OK | Successful read/retrieve operations or approved actions |
| **400** | Bad Request | Invalid request body or missing required fields |
| **403** | Forbidden | Missing required headers (e.g., X-User-Id) |
| **404** | Not Found | Resource not found |
| **500** | Internal Server Error | Database or service layer failures |

---

## Endpoints and Their Status Code Handling

### 1. **POST /api/v1/topics** - Create Topic

#### Request Header Validation
- **X-User-Id**: Required
  - ❌ Missing → **403 Forbidden** - `{"error":"X-User-Id header is required"}`

#### Request Body Validation
- Invalid JSON format
  - ❌ Invalid JSON → **400 Bad Request** - `{"error":"Invalid request body"}`

#### Field Validation
- `name`: Required, non-empty
  - ❌ Empty → **400 Bad Request** - `{"error":"Topic name is required"}`

- `cluster`: Required, non-empty
  - ❌ Empty → **400 Bad Request** - `{"error":"Cluster name is required"}`

- `partitions`: Required, must be > 0
  - ❌ <= 0 → **400 Bad Request** - `{"error":"Partitions must be greater than 0"}`

- `replicas`: Required, must be > 0
  - ❌ <= 0 → **400 Bad Request** - `{"error":"Replicas must be greater than 0"}`

#### Success Response
- ✅ All validations pass
  - Status: **201 Created**
  - Response: Topic object in JSON format

#### Service/DB Failures
- Database or service layer errors
  - Status: **500 Internal Server Error**
  - Response: `{"error":"Failed to create topic"}`

---

### 2. **GET /api/v1/topics** - List Topics

#### Success Response
- ✅ Topics retrieved successfully
  - Status: **200 OK**
  - Response: Array of topic objects (empty array if no topics)

#### Service/DB Failures
- Database or service layer errors
  - Status: **500 Internal Server Error**
  - Response: `{"error":"Failed to retrieve topics"}`

---

### 3. **GET /api/v1/topics/{name}** - Get Topic

#### Path Validation
- `name`: Required path parameter
  - ❌ Empty → **400 Bad Request** - `{"error":"Topic name is required"}`

#### Success Response
- ✅ Topic found
  - Status: **200 OK**
  - Response: Topic object in JSON format

#### Not Found
- Topic doesn't exist in database
  - Status: **404 Not Found**
  - Response: `{"error":"Topic not found"}`

---

### 4. **POST /api/v1/topics/{name}/approve** - Approve Topic

#### Path Validation
- `name`: Required path parameter
  - ❌ Empty → **400 Bad Request** - `{"error":"Topic name is required"}`

#### Request Header Validation
- **X-User-Id**: Required for approval
  - ❌ Missing → **403 Forbidden** - `{"error":"X-User-Id header is required"}`

#### Success Response
- ✅ Topic approved successfully
  - Status: **200 OK**
  - Response: `{"status":"approved"}`

#### Service/DB Failures
- Database or service layer errors
  - Status: **500 Internal Server Error**
  - Response: `{"error":"Failed to approve topic"}`

---

### 5. **POST /api/v1/policies** - Create Policy

#### Request Body Validation
- Invalid JSON format
  - ❌ Invalid JSON → **400 Bad Request** - `{"error":"Invalid request body"}`

#### Field Validation
- `principal`: Required, non-empty
  - ❌ Empty → **400 Bad Request** - `{"error":"Principal is required"}`

- `action`: Required, non-empty
  - ❌ Empty → **400 Bad Request** - `{"error":"Action is required"}`

- `resource`: Required, non-empty
  - ❌ Empty → **400 Bad Request** - `{"error":"Resource is required"}`

- `effect`: Required, must be "permit" or "forbid"
  - ❌ Other values → **400 Bad Request** - `{"error":"Effect must be either 'permit' or 'forbid'"}`

#### Success Response
- ✅ All validations pass
  - Status: **201 Created**
  - Response: Policy object in JSON format

#### Service/DB Failures
- Database or service layer errors
  - Status: **500 Internal Server Error**
  - Response: `{"error":"Failed to create policy"}`

---

## Response Format

### Success Response
```json
{
  "id": "123",
  "name": "orders",
  "cluster": "prod",
  "partitions": 3,
  "replicas": 2,
  "status": "PENDING",
  "requestedBy": "user123",
  "createdAt": "2025-12-29T22:19:06Z"
}
```

### Error Response
```json
{
  "error": "Error message here"
}
```

---

## Testing Examples

### Create Topic (Success)
```bash
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -H "X-User-Id: user123" \
  -d '{
    "name": "orders",
    "cluster": "prod",
    "partitions": 3,
    "replicas": 2
  }'
```

### Create Topic (Missing Header)
```bash
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "orders",
    "cluster": "prod",
    "partitions": 3,
    "replicas": 2
  }'
# Response: 403 Forbidden
```

### Create Topic (Invalid Data)
```bash
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -H "X-User-Id: user123" \
  -d '{
    "name": "",
    "cluster": "prod",
    "partitions": 0,
    "replicas": 2
  }'
# Response: 400 Bad Request (multiple validations)
```

### List Topics
```bash
curl -X GET http://localhost:8080/api/v1/topics
# Response: 200 OK with array of topics
```

### Get Topic
```bash
curl -X GET http://localhost:8080/api/v1/topics/orders
# Response: 200 OK (found) or 404 Not Found
```

### Approve Topic
```bash
curl -X POST http://localhost:8080/api/v1/topics/orders/approve \
  -H "X-User-Id: admin123"
# Response: 200 OK with status message
```

---

## Logging

Every endpoint logs:
- ✅ **INFO**: Entry and exit points of the API handler (middleware)
- ✅ **ERROR**: Validation errors and failures
- ✅ **DEBUG**: Successful parsing and validation steps

All logs include:
- Timestamp
- Log level (INFO, DEBUG, ERROR)
- File name and line number
- Log message

Color-coded output:
- 🟢 Green: INFO
- 🟡 Yellow: DEBUG
- 🔴 Red: ERROR
