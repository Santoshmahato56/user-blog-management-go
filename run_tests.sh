#!/bin/bash

# Define colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Installing dependencies...${NC}"
go get github.com/stretchr/testify/mock@v1.10.0
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/DATA-DOG/go-sqlmock

echo -e "${YELLOW}Running tests and generating coverage reports...${NC}"

# Create directory for coverage reports if it doesn't exist
mkdir -p coverage

# Run tests with coverage for each package
echo -e "${YELLOW}Running repository tests...${NC}"
go test -v -coverprofile=coverage/repository.out ./internal/repository/impl/...

echo -e "${YELLOW}Running service tests...${NC}"
go test -v -coverprofile=coverage/service.out ./internal/service/impl/...

# Merge coverage profiles
echo -e "${YELLOW}Merging coverage profiles...${NC}"
go test -v -coverprofile=coverage/coverage1.out ./internal/service/impl/... ./internal/repository/impl/...

# Generate HTML coverage report
echo -e "${YELLOW}Generating HTML coverage report...${NC}"
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Display total coverage
total_coverage=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}')
echo -e "${GREEN}Total test coverage: ${total_coverage}${NC}"

echo -e "${GREEN}Done! Coverage report is available at coverage/coverage.html${NC}" 