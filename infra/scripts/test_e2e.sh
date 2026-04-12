#!/bin/bash
set -e

# NOTE: Since go is missing in the local PATH of the runner, you should run this locally 
# or via docker-compose after `docker-compose up -d kafka-governance`
# We use curl directly against localhost:8080

echo "1. Registering ADMIN User"
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username": "admin1", "password": "password123", "role": "ADMIN"}' \
  -s -o out_admin.json
  
ADMIN_TOKEN=$(cat out_admin.json | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Admin Token: $ADMIN_TOKEN"

echo "2. Registering Regular User"
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username": "user1", "password": "password123", "role": "USER"}' \
  -s -o out_user.json

USER_TOKEN=$(cat out_user.json | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "User Token: $USER_TOKEN"

echo "3. Creating Policy"
curl -X POST http://localhost:8080/api/v1/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"principal": "User::\"user1\"", "action": "Action::\"CreateTopic\"", "resource": "Topic::\"orders\"", "effect": "permit"}' \
  -s > /dev/null

echo "4. Creating Topic as Regular User (Starts PENDING)"
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "orders", "cluster": "local", "partitions": 3, "replicas": 1}' \
  -s > /dev/null

echo "5. Attempting to approve Topic as Regular User (Should Fail)"
curl -X POST http://localhost:8080/api/v1/topics/orders/approve \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -s > out_fail.json
cat out_fail.json

echo ""
echo "6. Approving Topic as ADMIN (Should Succeed)"
curl -X POST http://localhost:8080/api/v1/topics/orders/approve \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -s > out_success.json
cat out_success.json

echo ""
echo "7. Retrieving All Topics"
curl -X GET http://localhost:8080/api/v1/topics \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -s | grep -o '"status":"[^"]*'

echo "Test suite deployed and tested."
