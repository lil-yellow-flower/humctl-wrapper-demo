package constants

// Config field names
const (
	HumanitecToken = "humanitec_token"
	HumanitecOrg   = "humanitec_org"
)

// Default values
const (
	DefaultOutputFormat = "table"
	DefaultConfigFile   = "$HOME/config.yaml"
)

// Command names
const (
	RootCmdName = "humctl-wrapper-demo"
	RootCmdUse  = "humctl-wrapper"

	// Get command
	GetCmdName = "get"
	GetCmdUse  = "get"

	// Get apps command
	GetAppsCmdUse = "apps"
)

// Command descriptions
const (
	GetCmdShort     = "Get resources from Humanitec platform"
	GetAppsCmdShort = "Get applications from Humanitec platform"
)

// Flag names
const (
	// Global flags
	ConfigFlagName   = "config"
	ConfigFlagShort  = "c"
	VersionFlagName  = "version"
	VersionFlagShort = "v"

	// Get apps flags
	OutputFlagName = "output"
	OrgFlagName    = "org"
)

// Flag shorthands
const (
	OutputFlagShort = "o"
	OrgFlagShort    = "g"
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
	ErrMissingToken        = "Humanitec API token is required"
	ErrMissingOrg          = "Humanitec organization ID is required"
	ErrInvalidOutputFormat = "invalid output format: %v"
	ErrInvalidOrgFlag      = "invalid organization flag: %v"
	ErrGetApps            = "failed to get applications: %v"
	ErrFormatOutput       = "failed to format output: %v"
	ErrClientInit         = "failed to initialize client: %v"
) 