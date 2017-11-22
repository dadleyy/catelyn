package catelyn

import "github.com/alecthomas/kingpin"

// Command implementations representan interface that can be used as a configurable kingpin action.
type Command interface {
	Exec(*kingpin.ParseContext) error
	Configure(*kingpin.CmdClause)
}
