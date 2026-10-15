#!/usr/bin/env bash

# vars
BASE_URL="http://localhost:8080"

# status check the response code and prints a proper response
# this assumes the final three characters are the HTTP status code.
status() {
    local req="$1"
    local url="$2"
    local body code

    body=$(curl -s -L -X "$req" -w "%{http_code}" "$url")
    code="${body: -3}" 
    body="${body::-3}"

    echo "Request: $req $url"
    echo "Response Code: $code"
    echo "Response Body:"
    echo "$body"
}

# happy paths
echo "========SUCCESSFUL REQUESTS========"
echo "GET /app"
status "GET" $BASE_URL/app

echo "GET /healthz"
status GET $BASE_URL/healthz

echo "GET /metrics"
status GET $BASE_URL/metrics

echo "POST /reset"
status POST $BASE_URL/reset

echo "GET /metrics (Hits: 0)"
status GET $BASE_URL/metrics

# error/sad paths
echo "========BAD REQUESTS========"

echo "POST /healthz (should be rejected)"
status "POST" "$BASE_URL/healthz"

echo "POST /healthz (should be rejected)"
status "POST" "$BASE_URL/metrics" 

echo "GET /reset (should be rejected)"
status "GET" "$BASE_URL/healthz" 
