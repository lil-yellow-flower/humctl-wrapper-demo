# Humctl Wrapper Demo

A command line interface wrapper for Humanitec platform that provides basic operations for managing resources.

## Prerequisites

- Go 1.21 or later
- A Humanitec account and API token

## Quick Start

```bash
# Clone the repository
git clone https://github.com/lil-yellow-flower/humctl-wrapper-demo.git
cd humctl-wrapper-demo

# Build the application
go build -o humctl-wrapper
```

## Create a .env file in root directory with your Humanitec credentials

```env
HUMANITEC_TOKEN=your-token-here
HUMANITEC_ORG=your-org-here
```

## Available Commands

### Get Applications

```bash
# List applications (uses HUMANITEC_ORG from .env)
./humctl-wrapper get apps

# List applications for a specific organization
./humctl-wrapper get apps --org different-org

# Output formats
./humctl-wrapper get apps -o json
./humctl-wrapper get apps -o yaml
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./internal/commands

# Run specific test
go test -v ./internal/commands -run TestGetApps/table_format
```

## License

MIT License - see [LICENSE](LICENSE)
