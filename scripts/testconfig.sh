#!/bin/bash
echo "Running configtest"

git diff --name-only HEAD~1 HEAD | grep config.yml &> /dev/null
if [[ $? -eq 0 ]]; then
  go run cmd/configtest/configtest.go -c config.yml
else
  echo "No config changes detected"
  exit 0
fi
