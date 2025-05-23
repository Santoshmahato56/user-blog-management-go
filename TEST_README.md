# User Blog Management Test Documentation

This document provides information about the testing strategy, structure, and how to run tests for the User Blog Management application.

## Test Structure

The tests follow a table-driven approach and are organized as follows:

- **Mock Definitions**: Located in `/internal/mocks/`, these files provide mock implementations of interfaces for testing.
- **Repository Tests**: Located in `/internal/repository/impl/`, these tests verify that repository implementations correctly interact with the database.
- **Service Tests**: Located in `/internal/service/impl/`, these tests verify that service implementations correctly handle business logic and interact with repositories.

## Running Tests

### Prerequisites

- Go 1.14 or higher
- All dependencies installed

### Running All Tests

To run all tests and generate a coverage report:

```bash
./run_tests.sh
```

This script will:
1. Run all repository tests
2. Run all service tests
3. Generate a merged coverage profile
4. Create an HTML coverage report
5. Display the total test coverage percentage

The HTML coverage report will be available at `coverage/coverage.html`

### Running Specific Tests

To run tests for a specific package:

```bash
go test -v ./internal/repository/impl/...
```

Or for a specific test file:

```bash
go test -v ./internal/repository/impl/blog_repository_test.go
```

### Test Coverage

To generate a coverage report for a specific package:

```bash
go test -coverprofile=coverage.out ./internal/service/impl/...
go tool cover -html=coverage.out
```

## Test Naming Convention

Tests follow a consistent naming convention:

```
Test<Package>_<Method>_<Scenario>
```

For example:
- `TestBlogService_Create_Success`
- `TestBlogRepository_FindByID_NotFound`

## Mock Usage

The tests use the `github.com/stretchr/testify/mock` package to create mock implementations of interfaces.

Example of setting up a mock:

```go
// Initialize mock repository
mockRepo := new(mocks.MockBlogRepository)

// Create the service with mock repository
blogService := impl.NewBlogService(mockRepo)

// Set expectations for mock calls
mockRepo.On("FindByID", uint(1)).Return(&models.Blog{}, nil)

// Call the service method that uses the mock
result, err := blogService.GetByID(uint(1))

// Assert that expectations were met
mockRepo.AssertExpectations(t)
```

## Database Mocking

The repository tests use `github.com/DATA-DOG/go-sqlmock` to mock the database layer, allowing tests to verify SQL queries without requiring an actual database connection.

Example:

```go
// Setup mock DB
mockDB, mock, err := setupMockDB(t)

// Set expectations for SQL queries
mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blogs`")).
    WillReturnRows(rows)

// Call repository method
blogs, count, err := repo.List(0, 10, true)

// Verify all expectations were met
mock.ExpectationsWereMet()
```