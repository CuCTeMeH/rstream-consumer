#!/bin/bash

# List all packages and exclude the ones you want to ignore
packages=$(go list ./... | grep -v 'github.com/cuctemeh/rstream-consumer/internal/testing/mocks' | grep -v 'github.com/cuctemeh/rstream-consumer/internal/shutdown' | grep -v 'github.com/cuctemeh/rstream-consumer/internal/storage' | grep -v 'github.com/cuctemeh/rstream-consumer/internal/config' | grep -v 'github.com/cuctemeh/rstream-consumer/cmd')

# Run tests and generate coverage profile, including only the specified packages
go test -coverprofile=cover.out -coverpkg=$(echo $packages | tr ' ' ',') ./...

# Extract the total coverage percentage
coverage=$(go tool cover -func=cover.out | grep total | awk '{print $3}')

# Check if coverage is empty
if [ -z "$coverage" ]; then
  echo "Failed to calculate coverage"
  exit 1
fi

# Print the total coverage
echo "Total coverage: $coverage"

# Convert coverage to a number for comparison
coverage_number=$(echo $coverage | sed 's/%//')

# Check if coverage is below the threshold
if (( $(echo "$coverage_number < 70.0" | bc -l) )); then
  echo "Coverage is below the threshold"
  exit 1
fi