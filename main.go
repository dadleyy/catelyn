package main

import "io"
import "os"
import "fmt"
import "github.com/bgentry/speakeasy"
import "github.com/alecthomas/kingpin"
import "github.com/dadleyy/catelyn/catelyn"
import "github.com/dadleyy/catelyn/catelyn/constants"

type commandConfiguration struct {
	name        string
	description string
	flags       uint
}

func main() {
	var out io.Writer = os.Stdout

	options := catelyn.GlobalCLIOptions{}

	cli := kingpin.New("catelyn", "A confluence helper")
	cli.Flag("username", "your confluence username").Short('u').Required().StringVar(&options.ConfluenceUsername)
	cli.Flag("hostname", "your confluence hostname").Short('h').Required().StringVar(&options.ConfluenceHost)
	cli.Flag("password", "your confluence password").Short('p').StringVar(&options.ConfluencePassword)

	loadPassword := func(context *kingpin.ParseContext) (e error) {
		if options.ConfluencePassword != "" {
			return nil
		}

		options.ConfluencePassword, e = speakeasy.FAsk(out, constants.PasswordPrompt)
		fmt.Fprintln(out)
		return
	}

	commands := map[commandConfiguration]catelyn.Command{
		commandConfiguration{
			name:        "search-spaces",
			description: "Search confluence spaces",
			flags:       constants.RequirePassword,
		}: catelyn.NewSearchSpacesCommand(out, &options),
		commandConfiguration{
			name:        "search-pages",
			description: "Search confluence pages",
			flags:       constants.RequirePassword,
		}: catelyn.NewSearchPagesCommand(out, &options),
	}

	for config, command := range commands {
		item := cli.Command(config.name, config.description)

		if config.flags&constants.RequirePassword != 0 {
			item.Action(loadPassword)
		}

		command.Configure(item)
		item.Action(command.Exec)
	}

	kingpin.MustParse(cli.Parse(os.Args[1:]))
}
