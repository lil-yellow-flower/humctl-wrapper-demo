package constants

// Environment variable names
const (
	// Humanitec environment variables
	EnvHumanitecToken = "HUMANITEC_TOKEN"
	EnvHumanitecOrg   = "HUMANITEC_ORG"
	EnvHumanitecEnv   = "HUMANITEC_ENV"
)

// Default values
const (
	DefaultOutputFormat = "table"
	DefaultConfigFile   = "$HOME/.humctl-wrapper.yaml"
)

// Command names and paths
const (
	// Root command
	RootCmdName = "humctl-wrapper-demo"
	RootCmdUse  = "humctl-wrapper"

	// Get command
	GetCmdName = "get"
	GetCmdUse  = "get"

	// Get apps command
	GetAppsCmdUse = "apps"
)

// Flag names
const (
	// Global flags
	ConfigFlagName   = "config"
	ConfigFlagShort  = "c"
	VersionFlagName  = "version"
	VersionFlagShort = "v"

	// Get apps flags
	OutputFlagName  = "output"
	OutputFlagShort = "o"
	OrgFlagName     = "org"
)

// Help text
const (
	// Global help text
	ConfigFlagHelp   = "config file (default is $HOME/.humctl-wrapper.yaml)"
	VersionFlagHelp  = "Print the version number"

	// Get apps help text
	OutputFlagHelp = "Output format (table|json|yaml)"
	OrgFlagHelp    = "Humanitec organization ID (defaults to %s environment variable)"
)

// Error messages
const (
	ErrMissingToken        = "HUMANITEC_TOKEN environment variable is required"
	ErrMissingOrg          = "HUMANITEC_ORG environment variable is required"
	ErrInvalidOutputFormat = "invalid output format: %v"
	ErrInvalidOrgFlag      = "invalid organization flag: %v"
	ErrGetApps            = "failed to get applications: %v"
	ErrFormatOutput       = "failed to format output: %v"
	ErrClientInit         = "failed to initialize client: %v"
) 