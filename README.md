# Humctl Wrapper Demo

A CLI wrapper for Humanitec platform that provides basic CRUD operations for managing resources.

## Prerequisites

- Go 1.21+
- Humanitec account and API token

## Setup

1. Clone and build:
```bash
git clone https://github.com/mathi-ma51zaw/humctl-wrapper-demo.git
cd humctl-wrapper-demo
go build -o humctl-wrapper.exe
```

2. Create `.env`:
```env
HUMANITEC_TOKEN=your_token_here
HUMANITEC_ORG=your_org_here
HUMANITEC_ENV=your_env_here
```

## Usage

```bash
# List applications
./humctl-wrapper.exe get apps

# Output formats
./humctl-wrapper.exe get apps -o json
./humctl-wrapper.exe get apps -o yaml
```

## Project Structure

```
humctl-wrapper-demo/
├── main.go            # Entry point
├── internal/
│   ├── commands/      # Command implementations
│   ├── humanitec/     # Humanitec client
│   └── output/        # Output formatting
├── go.mod
└── README.md
```

## License

MIT License - see [LICENSE](LICENSE)
