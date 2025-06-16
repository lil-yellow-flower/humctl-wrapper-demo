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

// Command use strings
const (
	RootCmdUse   = "humctl-wrapper"
	GetCmdUse    = "get"
	CreateCmdUse = "create"
	UpdateCmdUse = "update"
	DeleteCmdUse = "delete"
	AppsCmdUse   = "apps"
	AppCmdUse    = "app"
)

// Command short descriptions
const (
	RootCmdShort   = "A wrapper for the Humanitec CLI"
	GetCmdShort    = "Get resources"
	CreateCmdShort = "Create resources"
	UpdateCmdShort = "Update resources"
	DeleteCmdShort = "Delete resources"
	AppsCmdShort   = "Manage applications"
	AppCmdShort    = "Manage a single application"
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

	// Create app flags
	NameFlagName           = "name"
	SkipEnvCreationFlagName = "skip-env-creation"
	IDFlagName             = "id"

	// Flag shorthands
	OutputFlagShort = "o"
	OrgFlagShort    = "g"
	NameFlagShort   = "n"
	SkipEnvCreationFlagShort = "s"
	IDFlagShort             = "i"
)

// Help text
const (
	// Global help text
	ConfigFlagHelp   = "config file (default is $HOME/.humctl-wrapper.yaml)"
	VersionFlagHelp  = "Print the version number"

	// Get apps help text
	OutputFlagHelp = "Output format (table|json|yaml)"
	OrgFlagHelp    = "Humanitec organization ID (defaults to %s environment variable)"

	// Create app help text
	NameFlagHelp           = "Name of the application"
	SkipEnvCreationFlagHelp = "Skip environment creation"
	IDFlagHelp             = "Application ID"
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
	ErrCreateApp          = "failed to create application: %v"
	ErrDeleteApp          = "failed to delete application: %v"
	ErrUpdateApp          = "failed to update application: %v"
	ErrLoadConfig         = "failed to load config: %v"
	ErrMissingAPIToken    = "API token is required"
	ErrAPIError           = "API error"
	ErrGetApp             = "failed to get application: %v"

	// Create command error messages
	CreateErrorMissingID   = "Error: required flag(s) \"id\" not set"
	CreateErrorMissingName = "Error: required flag(s) \"name\" not set"
	CreateErrorInvalidID   = "Error: invalid id format: must match pattern ^[a-z0-9](?:-?[a-z0-9]+)+$"
	CreateErrorIDTooLong   = "Error: invalid id format: must be 50 characters or less"
	CreateErrorDuplicateID = "Error: application with id '%s' already exists"
	CreateErrorServerError = "Error: failed to create app: %s"
)

// Success messages
const (
	SuccessAppUpdated = "Application successfully updated"
)

// Config-related constants
const (
	ConfigDir = "humctl"
	ConfigFile = "config.yaml"
) 