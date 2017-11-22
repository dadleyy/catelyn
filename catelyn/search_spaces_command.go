package catelyn

import "io"
import "fmt"
import "strings"
import "net/url"
import "text/tabwriter"
import "github.com/alecthomas/kingpin"

// NewSearchSpacesCommand returns an executable kingpin action for searching spaces.
func NewSearchSpacesCommand(out io.Writer, config *GlobalCLIOptions) Command {
	cli := searchSpacesCLI{
		globals: config,
		output:  out,
	}

	return &cli
}

type searchSpacesCLI struct {
	globals *GlobalCLIOptions
	output  io.Writer
	search  ConfluenceSpaceSearchInput
}

func (c *searchSpacesCLI) Configure(clause *kingpin.CmdClause) {
	clause.Flag("type", "the type of spaces to search ('personal', 'global')").Short('t').StringVar(&c.search.Type)
	clause.Flag("limit", "how many results to return").Default("10").Short('l').Uint8Var(&c.search.Limit)
	clause.Flag("start", "how many results to skip").Short('o').Uint8Var(&c.search.Start)
}

func (c *searchSpacesCLI) Exec(context *kingpin.ParseContext) error {
	uinfo := url.UserPassword(c.globals.ConfluenceUsername, c.globals.ConfluencePassword)
	client, e := NewConfluenceClient(uinfo, c.globals.ConfluenceHost)

	if e != nil {
		return e
	}

	spaces, paging, e := client.SearchSpaces(&c.search)
	writer := tabwriter.NewWriter(c.output, 10, 2, 3, ' ', 0)

	if e != nil {
		return e
	}

	fmt.Fprintf(c.output, "Spaces found on %s:\n", c.globals.ConfluenceHost)
	fmt.Fprintf(writer, "Space Key\tSpace Name\n")

	for _, s := range spaces {
		line := strings.Join([]string{s.Key, s.Name}, "\t")
		fmt.Fprintf(writer, "%s\n", line)
	}

	writer.Flush()

	fmt.Fprintf(c.output, "more results: %t\n", len(paging.Links.Next) > 0)

	return nil
}
