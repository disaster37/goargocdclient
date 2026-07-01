# Contributing

## Prerequisites

- Go 1.26+
- [golangci-lint](https://golangci-lint.run/) (optional, for advanced linting)

## Project Structure

```
.
├── client.go              # Main client constructor and options
├── client_test.go         # Client tests
├── api/
│   ├── api.go             # API interface and implementation
│   ├── account.go         # Account sub-client
│   ├── application.go     # Application sub-client
│   ├── applicationset.go  # ApplicationSet sub-client
│   ├── certificate.go     # Certificate sub-client
│   ├── cluster.go         # Cluster sub-client
│   ├── errors.go          # Error types and helpers
│   ├── gpgkey.go          # GPGKey sub-client
│   ├── models.go          # Shared model types
│   ├── notification.go    # Notification sub-client
│   ├── project.go         # Project sub-client
│   ├── repository.go      # Repository sub-client
│   ├── repocreds.go       # RepoCreds sub-client
│   ├── session.go         # Session sub-client
│   ├── settings.go        # Settings sub-client
│   ├── version.go         # Version sub-client
│   └── *_test.go          # Tests for each sub-client
└── go.mod
```

## Development Commands

### Build

```bash
go build ./...
```

### Test

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests for a specific package:

```bash
go test -v ./api/
```

### Code Coverage

Generate a coverage profile:

```bash
go test -coverprofile=coverage.out ./...
```

View coverage summary per package:

```bash
go test -cover ./...
```

View coverage in the terminal (by function):

```bash
go tool cover -func=coverage.out
```

Generate an HTML coverage report:

```bash
go tool cover -html=coverage.out -o coverage.html
```

### Lint

Run Go vet (static analysis):

```bash
go vet ./...
```

Run golangci-lint (if installed):

```bash
golangci-lint run ./...
```

### Format Code

```bash
go fmt ./...
```

### All-in-One Check

Run the full pipeline before committing:

```bash
go fmt ./... && go vet ./... && go build ./... && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

## Rules

- Use `resty.Client` for all API calls
- All sub-clients must be defined as interfaces
- 100% code coverage is required — write tests for every code path
- Document usage with examples
- Use `fmt.Sprintf` instead of string concatenation
- All files in `api/` are dedicated to a single API element (e.g. `application.go` for the Application API)

## Pre-PR Checklist

- [ ] Code follows Golang best practices
- [ ] `go fmt ./...` passes
- [ ] `go vet ./...` passes
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] Code coverage is 100%
- [ ] Documentation is up to date with usage examples
- [ ] No secrets or credentials are committed
- [ ] Error handling covers all failure paths

## How to Contribute

1. Fork this repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes following the project rules
4. Run the full pipeline (`go fmt && go vet && go build && go test -cover ./...`)
5. Commit your changes
6. Push and open a Pull Request
