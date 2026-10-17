#!/usr/bin/env bash

# --- Vars ---
BASE_URL_APP="http://localhost:8080"
BASE_URL_API="http://localhost:8080/api"
BASE_URL_ADMIN="http://localhost:8080/admin"


# --- Functions ---
# -- Status Check --
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
# -- Run tests --
# runTest orchestrates the test runs 
runTest() {
    local -n cases="$1"
    local -n inputs="$2"
    for ((i = 0; i < ${#cases[@]}; i++)); do
        echo "${cases[$i]}"
        method_index=$(( i * 2 ))
        url_index=$(( i * 2 + 1 ))
        status "${inputs[$method_index]}" "${inputs[$url_index]}"
        echo "========================"
    done
}

# --- Test cases ---
# -- Happy path --
test_case_happy=(
    "POST /admin/reset"
    "GET /admin/metrics (FileserverRequests: 0)"
    "GET /app"
    "GET /admin/metrics (FileserverRequests: 1)"
    "GET /app/assets/"
    "GET /admin/metrics (FileserverRequests: 2)"
    "GET /api/healthz"
    "GET /admin/metrics (FileserverRequests: 2)" 
)
test_input_happy=(
    "POST" "$BASE_URL_ADMIN/reset"
    "GET" "$BASE_URL_ADMIN/metrics"
    "GET" "$BASE_URL_APP/app"
    "GET" "$BASE_URL_ADMIN/metrics"
    "GET" "$BASE_URL_APP/app/assets/"
    "GET" "$BASE_URL_ADMIN/metrics"
    "GET" "$BASE_URL_API/healthz"
    "GET" "$BASE_URL_ADMIN/metrics"
)
# -- Error/sad path --
test_case_bad=(
    "POST /api/healthz (should be rejected)"
    "POST /api/healthz (should be rejected)"
    "GET /admin/reset (should be rejected)"
) 
test_input_bad=(
    "POST" "$BASE_URL_API/healthz"
    "POST" "$BASE_URL_ADMIN/metrics"
    "GET" "$BASE_URL_ADMIN/reset"  
) 

# --- Run tests ---
# -- Happy path --
echo "========SUCCESSFUL REQUESTS========"
runTest "test_case_happy" "test_input_happy"
# -- Error/sad path --
echo "========BAD REQUESTS========"
runTest "test_case_bad" "test_input_bad"