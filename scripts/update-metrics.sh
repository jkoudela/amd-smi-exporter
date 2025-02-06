#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
else
    echo "Error: .env file not found"
    echo "Please copy .env.example to .env and update the values"
    exit 1
fi

# Check required variables
if [ -z "$GRAFANA_URL" ] || [ -z "$GRAFANA_API_KEY" ]; then
    echo "Error: GRAFANA_URL and GRAFANA_API_KEY must be set in .env file"
    exit 1
fi

# Validate JSON file
if ! jq empty grafana/amd_gpu_metrics.json 2>/dev/null; then
    echo "Error: Invalid JSON in dashboard file"
    exit 1
fi

# Prepare dashboard JSON with wrapper
echo "Preparing dashboard JSON..."
TMP_JSON=$(mktemp)
trap "rm -f $TMP_JSON" EXIT

jq -n --slurpfile dashboard grafana/amd_gpu_metrics.json \
  '{dashboard: $dashboard[0], folderId: 0, overwrite: true}' > "$TMP_JSON"

# Update dashboard
echo "Updating dashboard..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Authorization: Bearer $GRAFANA_API_KEY" \
    -H "Content-Type: application/json" \
    -d "@$TMP_JSON" \
    "${GRAFANA_URL}/api/dashboards/db")

# Get status code
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | sed \$d)

# Check response
if [ "$HTTP_CODE" -eq 200 ]; then
    echo "Dashboard updated successfully!"
    echo "Response: $BODY"
else
    echo "Error updating dashboard"
    echo "Status code: $HTTP_CODE"
    echo "Response: $BODY"
    exit 1
fi
