package constants

const (
	// PasswordPrompt is displayed before running any command.
	PasswordPrompt = "password: "
)

const (
	// RequirePassword lets the main application know that the command will require a user's password to run.
	RequirePassword = iota + 1

	// AllFlags lets the main application know that all other flags should be honored.
	AllFlags
)
