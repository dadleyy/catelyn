package main

import "io"
import "os"
import "fmt"
import "github.com/joho/godotenv"
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
	cli.Flag("username", "confluence username").Envar(constants.ConfluenceUsernameEnvironmentVariable).
		Short('u').Required().StringVar(&options.ConfluenceUsername)

	cli.Flag("hostname", "confluence hostname").Envar(constants.ConfluenceHostnameEnvironmentVariable).
		Short('h').Required().StringVar(&options.ConfluenceHost)

	cli.Flag("password", "confluence password").Envar(constants.ConfluencePasswordEnvironmentVariable).
		Short('p').StringVar(&options.ConfluencePassword)

	godotenv.Load()

	loadPassword := func(context *kingpin.ParseContext) (e error) {
		if options.ConfluencePassword != "" {
			return nil
		}

		fmt.Fprintf(out, "%s requires a password to continue\n", context.SelectedCommand.FullCommand())
		options.ConfluencePassword, e = speakeasy.Ask(constants.PasswordPrompt)
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
		commandConfiguration{
			name:        "create-page",
			description: "Creates a confluence page from a file",
			flags:       constants.RequirePassword,
		}: catelyn.NewCreatePageCommand(out, os.Stdin, &options),
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
