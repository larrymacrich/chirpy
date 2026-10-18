#!/usr/bin/env bash

# --- Vars ---
BASE_URL_APP="http://localhost:8080"
BASE_URL_API="http://localhost:8080/api"
BASE_URL_ADMIN="http://localhost:8080/admin"


# --- Functions ---
# -- Status Check --
# status check the response code and prints a proper response
# this assumes the final three characters are the HTTP status code.
# checking for GET vs POST requests for additional request data
#       add json data to POST Requests if the the body is non-empty (!= "")
status() {
    local req="$1"
    local url="$2"
    local data="$3"
    local body code curl_args=()

    curl_args=(-s -L -X "$req" -w "%{http_code}")   # default
    # add jsson data
    if [[ -n "$data" ]]; then
        curl_args+=(-H "Content-Type: application/json" -d "$data")
    fi

    body=$(curl "${curl_args[@]}" "$url")
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
        method_index=$(( i * 3 ))
        url_index=$(( i * 3 + 1 ))
        reqBody_index=$(( i * 3 + 2 ))
        status "${inputs[$method_index]}" "${inputs[$url_index]}" "${inputs[$reqBody_index]}"
        echo "========================"
    done
}

# --- Test cases ---
# Each test case = 3 consecutive elements: method, url, body (use "" if no body)
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
    'POST /api/validate_chirp {"body": "This is a test chirp"}'
    'POST /api/validate_chirp {"body": "I hear Mastodon is better than Chirpy. sharbert I need to migrate","extra": "this should be ignored"}'
    'POST /api/validate_chirp {"body": "I really need a kerfuffle to go to bed sooner, Fornax !"}'
)
test_input_happy=(
    "POST" "$BASE_URL_ADMIN/reset" ""
    "GET" "$BASE_URL_ADMIN/metrics" ""
    "GET" "$BASE_URL_APP/app" ""
    "GET" "$BASE_URL_ADMIN/metrics" ""
    "GET" "$BASE_URL_APP/app/assets/" ""
    "GET" "$BASE_URL_ADMIN/metrics" ""
    "GET" "$BASE_URL_API/healthz"  ""
    "GET" "$BASE_URL_ADMIN/metrics" ""
    "POST" "$BASE_URL_API/validate_chirp" '{"body": "This is a test chirp"}'
    "POST" "$BASE_URL_API/validate_chirp" '{"body": "I hear Mastodon is better than Chirpy. sharbert I need to migrate","extra": "this should be ignored"}'
    "POST" "$BASE_URL_API/validate_chirp" '{"body": "I really need a kerfuffle to go to bed sooner, Fornax !"}'
)
# -- Error/sad path --
test_case_bad=(
    "POST /api/healthz (should be rejected)"
    "POST /api/healthz (should be rejected)"
    "GET /admin/reset (should be rejected)"
    'POST/api/validate_chirp {"body": "lorem ipsum..."} (should be Status Code 400)'
) 
test_input_bad=(
    "POST" "$BASE_URL_API/healthz" ""
    "POST" "$BASE_URL_ADMIN/metrics" ""
    "GET" "$BASE_URL_ADMIN/reset"  ""
    "POST" "$BASE_URL_API/validate_chirp" '{"body": "lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."}'
)  

# --- Run tests ---
# -- Happy path --
echo "========SUCCESSFUL REQUESTS========"
runTest "test_case_happy" "test_input_happy"
# -- Error/sad path --
echo "========BAD REQUESTS========"
runTest "test_case_bad" "test_input_bad"