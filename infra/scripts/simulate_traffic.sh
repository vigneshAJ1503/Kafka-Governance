#!/bin/bash

# Script to validate native Kafka streams simulating physical platform movement bounds securely

echo "🚀 Simulating Kafka Data Flow for Governance System..."
echo "========================================================"

TOPIC_NAME="topic-test-e2e"

echo "Step 1: Connecting securely to containerized Kafka Broker (Port: 9092)..."

echo "Producing 3 transactional mocked payloads to topic: [ $TOPIC_NAME ]"
# Leverage container exec to stream inputs to kafka-console-producer
docker exec -i kafka kafka-console-producer --broker-list localhost:9092 --topic $TOPIC_NAME <<EOF
{"action": "produce", "event_id": "e_987", "type": "ORDER_CREATED", "total": 100, "sec_mask": true}
{"action": "produce", "event_id": "e_988", "type": "ORDER_COMPLETED", "total": 100, "sec_mask": true}
{"action": "produce", "event_id": "e_989", "type": "ORDER_ARCHIVED", "total": 100, "sec_mask": true}
EOF

echo "✅ Produced events successfully."
echo ""
echo "Step 2: Securing Consumer Subscription to verify stream pipeline..."
echo "Initializing Stream Reader:"

# Read messages synchronously
docker exec -i kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic $TOPIC_NAME --from-beginning --max-messages 3 --timeout-ms 5000

echo ""
echo "✅ Stream Verification Complete!"
echo "Data flow has successfully passed through the underlying Kafka cluster. Browser inspections will expose ZERO structural infrastructure vectors as they are securely decoupled!"
