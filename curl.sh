#!/usr/bin/env bash



curl -X POST \
  http://localhost:8081/accounts \
  -H 'Content-Type: application/json' \
  -d '{
  "owner": "jorge", "currency": "EUR"
}'