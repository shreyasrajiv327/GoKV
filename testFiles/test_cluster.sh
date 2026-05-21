#!/bin/bash

echo "=============================="
echo "GoKV Cluster Test"
echo "=============================="

BASE="http://localhost:8001"

echo ""
echo "1. Node Status"
curl -s "$BASE/node/status"
echo ""

echo ""
echo "2. PUT key=user value=100"
curl -s \
-X PUT "$BASE/kv" \
-H "Content-Type: application/json" \
-d '{
"key":"user",
"value":"100"
}'
echo ""

sleep 1

echo ""
echo "3. GET key=user"
curl -s "$BASE/kv/user"
echo ""

echo ""
echo "4. UPDATE key=user"
curl -s \
-X PUT "$BASE/kv" \
-H "Content-Type: application/json" \
-d '{
"key":"user",
"value":"200"
}'
echo ""

sleep 1

echo ""
echo "5. VERIFY UPDATE"
curl -s "$BASE/kv/user"
echo ""

echo ""
echo "6. DELETE key=user"
curl -s \
-X DELETE "$BASE/kv/user"
echo ""

sleep 1

echo ""
echo "7. VERIFY DELETE"
curl -s "$BASE/kv/user"
echo ""

echo ""
echo "8. WAL FILE"
ls -lh gokv.wal

echo ""
echo "9. SNAPSHOT FILE"
ls -lh snapshot.json

echo ""
echo "=============================="
echo "TEST COMPLETE"
echo "=============================="