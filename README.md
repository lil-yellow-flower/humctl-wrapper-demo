# Humctl Wrapper Demo

A command line interface wrapper for the Humanitec platform, demonstrating best practices for Go CLI development.

## Prerequisites

- Go 1.21 or later
- A Humanitec account with API access
- Your Humanitec API token and organization ID

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/lil-yellow-flower/humctl-wrapper-demo.git
   cd humctl-wrapper-demo
   ```

## Building

Build the application:
```bash
go build -o humctl-wrapper
```

## Usage

The CLI provides commands to interact with the Humanitec platform. Here are the different ways to use the `get apps` command:

```bash
# Get all applications (uses org from config.yaml)
./humctl-wrapper get apps

# Get applications for a specific organization
./humctl-wrapper get apps --org your-org-id

# Get applications in different output formats
./humctl-wrapper get apps --output table  # Default format
./humctl-wrapper get apps --output json   # JSON format
./humctl-wrapper get apps --output yaml   # YAML format

# Combine options
./humctl-wrapper get apps --org your-org-id --output json
```

## Configuration

The CLI can be configured using a `config.yaml` file in the project root:

```yaml
# Humanitec API credentials
humanitec_token: "your-api-token-here"
humanitec_org: "your-org-id-here"

# Default output format (table, json, yaml)
default_output: "table"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
