package constants

const (
	// PasswordPrompt is displayed before running any command.
	PasswordPrompt = "password: "

	// ConfluenceUsernameEnvironmentVariable is the environment variable name used for the confluence username.
	ConfluenceUsernameEnvironmentVariable = "CONFLUENCE_USERNAME"

	// ConfluencePasswordEnvironmentVariable is the environment variable name used for the confluence password.
	ConfluencePasswordEnvironmentVariable = "CONFLUENCE_PASSWORD"

	// ConfluenceHostnameEnvironmentVariable is the environment variable name used for the confluence hostname.
	ConfluenceHostnameEnvironmentVariable = "CONFLUENCE_HOSTNAME"
)

const (
	// RequirePassword lets the main application know that the command will require a user's password to run.
	RequirePassword = iota + 1

	// AllFlags lets the main application know that all other flags should be honored.
	AllFlags
)
