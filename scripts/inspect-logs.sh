#!/usr/bin/env bash

set -u

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
RESET="\033[0m"

SERVICES=(
  auth-service
  appointment-service
  doctor-service
  patient-service
  chat-service
  postgres
)

TAIL_LINES="${TAIL_LINES:-500}"
SINCE="${SINCE:-10m}"
RESTART_THRESHOLD="${RESTART_THRESHOLD:-3}"
REPORT_FILE="${REPORT_FILE:-log_inspection_report.txt}"
FAILED=0

print_ok() {
  printf "%b✓%b %s\n" "$GREEN" "$RESET" "$1"
}

print_warn() {
  printf "%b!%b %s\n" "$YELLOW" "$RESET" "$1"
}

print_fail() {
  printf "%b✗%b %s\n" "$RED" "$RESET" "$1"
  FAILED=1
}

write_report_header() {
  {
    echo "Log-Based Troubleshooting Automation Report"
    echo "Generated at: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
    echo "Log time window inspected per service: $SINCE"
    echo "Log tail lines inspected per service: $TAIL_LINES"
    echo
  } > "$REPORT_FILE"
}

append_report() {
  echo "$*" >> "$REPORT_FILE"
}

check_docker_compose() {
  if docker compose version >/dev/null 2>&1; then
    print_ok "Docker Compose is available"
  else
    print_fail "Docker Compose is not available"
    exit 1
  fi
}

inspect_service_logs() {
  local service="$1"
  local logs

  append_report "Service: $service"
  append_report "----------------------------------------"

  if ! logs="$(docker compose logs --no-color --since "$SINCE" --tail "$TAIL_LINES" "$service" 2>&1)"; then
    print_fail "$service logs could not be collected"
    append_report "LOG COLLECTION FAILED"
    append_report "$logs"
    append_report
    return
  fi

  if echo "$logs" | grep -Eiq "connection refused|could not connect|database unavailable|failed to setup|dial tcp|no such host|pq:|password authentication failed"; then
    print_fail "$service has possible database connection failures"
    append_report "Database connection failure patterns found:"
    echo "$logs" | grep -Ein "connection refused|could not connect|database unavailable|failed to setup|dial tcp|no such host|pq:|password authentication failed" | tail -20 >> "$REPORT_FILE"
  else
    print_ok "$service has no database connection failure patterns"
    append_report "No database connection failure patterns found."
  fi

  if echo "$logs" | grep -Eiq "panic|fatal|failed to run|container started|restarting|exited|segmentation violation|out of memory|oom"; then
    print_warn "$service has possible crash or restart-loop patterns"
    append_report "Crash/restart-loop patterns found:"
    echo "$logs" | grep -Ein "panic|fatal|failed to run|container started|restarting|exited|segmentation violation|out of memory|oom" | tail -20 >> "$REPORT_FILE"
  else
    print_ok "$service has no crash or restart-loop log patterns"
    append_report "No crash or restart-loop log patterns found."
  fi

  append_report
}

inspect_restart_count() {
  local service="$1"
  local container_id restart_count

  container_id="$(docker compose ps -q "$service" 2>/dev/null || true)"
  if [[ -z "$container_id" ]]; then
    print_warn "$service container is not running or was not found"
    append_report "$service restart count: unavailable"
    return
  fi

  restart_count="$(docker inspect --format '{{.RestartCount}}' "$container_id" 2>/dev/null || echo "unknown")"
  if [[ "$restart_count" =~ ^[0-9]+$ ]] && (( restart_count > RESTART_THRESHOLD )); then
    print_fail "$service restarted $restart_count times"
  else
    print_ok "$service restart count is $restart_count"
  fi

  append_report "$service restart count: $restart_count"
}

echo "Inspecting Docker Compose logs for troubleshooting patterns..."
check_docker_compose
write_report_header

for service in "${SERVICES[@]}"; do
  echo
  echo "Checking $service..."
  inspect_restart_count "$service"
  inspect_service_logs "$service"
done

echo
echo "Report saved to $REPORT_FILE"

if [[ "$FAILED" -ne 0 ]]; then
  print_fail "Log inspection found issues that need investigation"
  exit 1
fi

print_ok "Log inspection completed without critical issues"
