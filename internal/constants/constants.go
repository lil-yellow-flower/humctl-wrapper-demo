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
	RootCmdUse  = "humctl"

	// Get command
	GetCmdName = "get"
	GetCmdUse   = "get"

	// Get apps command
	GetAppsCmdUse = "apps"

	// Add command
	AddCmdName = "add"
	AddCmdUse   = "add"

	// Add app command
	AddAppCmdUse = "app"
)

// Command descriptions
const (
	RootCmdShort    = "A command line interface wrapper for Humanitec platform"
	GetCmdShort     = "Get resources from Humanitec platform"
	GetAppsCmdShort = "Get applications from Humanitec platform"
	AddCmdShort     = "Add resources to Humanitec platform"
	AddAppCmdShort  = "Add application to Humanitec platform"
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

	// Add app flags
	NameFlagName           = "name"
	SkipEnvCreationFlagName = "skip-env-creation"
)

// Flag shorthands
const (
	OutputFlagShort = "o"
	OrgFlagShort    = "g"
	NameFlagShort   = "n"
	SkipEnvCreationFlagShort = "s"
)

// Help text
const (
	// Global help text
	ConfigFlagHelp   = "config file (default is $HOME/.humctl-wrapper.yaml)"
	VersionFlagHelp  = "Print the version number"

	// Get apps help text
	OutputFlagHelp = "Output format (table|json|yaml)"
	OrgFlagHelp    = "Humanitec organization ID (defaults to %s environment variable)"

	// Add app help text
	NameFlagHelp           = "Name of the resource"
	SkipEnvCreationFlagHelp = "Skip environment creation"
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
	ErrInvalidName        = "invalid name: %v"
	ErrInvalidSkipEnvCreation = "invalid skip-env-creation flag: %v"
	ErrAddApp             = "failed to add application: %v"
	ErrLoadConfig         = "failed to load config: %v"
)

// Config-related constants
const (
	ConfigDir = "humctl"
	ConfigFile = "config.json"
) 