#!/usr/bin/env bash

set -u

GREEN="\033[0;32m"
RED="\033[0;31m"
RESET="\033[0m"

ENV_FILE=".env"
EXAMPLE_FILE=".env.example"
FAILED=0

declare -A ENV_VALUES

trim() {
  local value="$*"
  value="${value#"${value%%[![:space:]]*}"}"
  value="${value%"${value##*[![:space:]]}"}"
  printf '%s' "$value"
}

print_ok() {
  printf "%b✓%b %s\n" "$GREEN" "$RESET" "$1"
}

print_fail() {
  printf "%b✗%b %s\n" "$RED" "$RESET" "$1"
  FAILED=1
}

load_env() {
  local line key value

  while IFS= read -r line || [[ -n "$line" ]]; do
    line="$(trim "$line")"
    [[ -z "$line" || "$line" == \#* ]] && continue

    line="${line#export }"
    [[ "$line" != *=* ]] && continue

    key="$(trim "${line%%=*}")"
    value="$(trim "${line#*=}")"
    value="${value%\"}"
    value="${value#\"}"
    value="${value%\'}"
    value="${value#\'}"

    [[ -n "$key" ]] && ENV_VALUES["$key"]="$value"
  done < "$ENV_FILE"
}

key_is_defined() {
  local key="$1"
  [[ -n "${ENV_VALUES[$key]+set}" && -n "$(trim "${ENV_VALUES[$key]}")" ]]
}

check_key() {
  local key="$1"

  if key_is_defined "$key"; then
    print_ok "$key is defined"
  else
    print_fail "$key is missing or empty"
  fi
}

check_example_keys() {
  local line key

  while IFS= read -r line || [[ -n "$line" ]]; do
    line="$(trim "$line")"
    [[ -z "$line" || "$line" == \#* ]] && continue

    line="${line#export }"
    [[ "$line" != *=* ]] && continue

    key="$(trim "${line%%=*}")"
    [[ -n "$key" ]] && check_key "$key"
  done < "$EXAMPLE_FILE"
}

check_database_connection() {
  local host="${ENV_VALUES[DB_HOST]:-}"
  local port="${ENV_VALUES[DB_PORT]:-}"

  if [[ -z "$host" || -z "$port" ]]; then
    print_fail "Database TCP check skipped because DB_HOST or DB_PORT is missing"
    return
  fi

  if timeout 5 bash -c 'cat < /dev/null > /dev/tcp/$1/$2' _ "$host" "$port" 2>/dev/null; then
    print_ok "TCP connection to $host:$port succeeded"
  else
    print_fail "TCP connection to $host:$port failed"
  fi
}

echo "Validating healthcare microservices configuration..."

if [[ -s "$ENV_FILE" ]]; then
  print_ok "$ENV_FILE exists and is not empty"
else
  print_fail "$ENV_FILE is missing or empty"
fi

if [[ -s "$EXAMPLE_FILE" ]]; then
  print_ok "$EXAMPLE_FILE exists and is not empty"
else
  print_fail "$EXAMPLE_FILE is missing or empty"
fi

if [[ -s "$ENV_FILE" && -s "$EXAMPLE_FILE" ]]; then
  load_env

  echo
  echo "Checking keys from $EXAMPLE_FILE..."
  check_example_keys

  echo
  echo "Checking critical runtime variables..."
  CRITICAL_KEYS=(
    DB_HOST
    DB_PORT
    DB_NAME
    DB_USER
    DB_PASSWORD
    AUTH_SERVICE_URL
    APPOINTMENT_SERVICE_URL
    DOCTOR_SERVICE_URL
    PATIENT_SERVICE_URL
    CHAT_SERVICE_URL
  )

  for key in "${CRITICAL_KEYS[@]}"; do
    check_key "$key"
  done

  echo
  echo "Checking database connectivity..."
  check_database_connection
fi

echo
echo "Note: this pre-deployment validation prevents the same class of misconfiguration incident documented in Assignment 4 by catching missing service URLs and database settings before containers are started."

if [[ "$FAILED" -ne 0 ]]; then
  echo
  print_fail "Configuration validation failed"
  exit 1
fi

echo
print_ok "Configuration validation passed"
