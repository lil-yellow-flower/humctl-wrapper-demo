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
   ```
   ```bash
   cd humctl-wrapper-demo
   ```

## Building

Build the application:
```bash
go build -o humctl-wrapper.exe
```

## Configuration

The CLI requires a `config.yaml` file in the project root to run. Create this file with the following content:

```yaml
# Humanitec API credentials
humanitec_token: "your-api-token-here"
humanitec_org: "your-org-id-here"

# Default output format (table, json, yaml)
default_output: "table"
```

## Usage

The CLI provides commands to interact with the Humanitec platform. All commands support the following output formats:
- `--output table` (default)
- `--output json`
- `--output yaml`

### Get Applications

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

### Add Application

```bash
# Add a new application (uses org from config.yaml)
./humctl-wrapper add app --name "My Application"

# Add application for a specific organization
./humctl-wrapper add app --name "My Application" --org your-org-id

# Add application without creating default environment
./humctl-wrapper add app --name "My Application" --skip-env-creation

# Add application with different output formats
./humctl-wrapper add app --name "My Application" --output json
./humctl-wrapper add app --name "My Application" --output yaml

# Using shorthand flags
./humctl-wrapper add app -n "My Application" -g your-org-id -s -o json
```

### Delete Application

```bash
# Delete an application (uses org from config.yaml)
./humctl-wrapper delete app --name my-app-id

# Delete application for a specific organization
./humctl-wrapper delete app --name my-app-id --org your-org-id

# Delete application with different output formats
./humctl-wrapper delete app --name my-app-id --output json
./humctl-wrapper delete app --name my-app-id --output yaml

# Using shorthand flags
./humctl-wrapper delete app -n my-app-id -g your-org-id -o json
```

### Update Application

```bash
# Update an application name (uses org from config.yaml)
./humctl-wrapper update app --name my-app-id --new-name "Updated App Name"

# Update application for a specific organization
./humctl-wrapper update app --name my-app-id --new-name "Updated App Name" --org your-org-id

# Update application with different output formats
./humctl-wrapper update app --name my-app-id --new-name "Updated App Name" --output json
./humctl-wrapper update app --name my-app-id --new-name "Updated App Name" --output yaml

# Using shorthand flags
./humctl-wrapper update app -n my-app-id -m "Updated App Name" -g your-org-id -o json
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
