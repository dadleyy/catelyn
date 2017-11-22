package catelyn

import "io"
import "log"
import "fmt"
import "strings"
import "net/url"
import "text/tabwriter"
import "github.com/alecthomas/kingpin"

// NewSearchSpacesCommand returns an executable kingpin action.
func NewSearchSpacesCommand(out io.Writer, config *GlobalCLIOptions) Command {
	logger := log.New(out, "", 0)

	cli := searchCLI{
		globals: config,
		Logger:  logger,
		output:  out,
	}

	return &cli
}

type searchCLI struct {
	*log.Logger
	globals *GlobalCLIOptions
	query   string
	output  io.Writer
}

func (c *searchCLI) Configure(clause *kingpin.CmdClause) {
	f := clause.Flag("query", "a search term used to query spaces")
	f.Short('q')
	f.StringVar(&c.query)
}

func (c *searchCLI) Exec(context *kingpin.ParseContext) error {
	uinfo := url.UserPassword(c.globals.ConfluenceUsername, c.globals.ConfluencePassword)
	client, e := NewConfluenceClient(uinfo, c.globals.ConfluenceHost)

	if e != nil {
		return e
	}

	spaces, paging, e := client.SearchSpaces(c.query)
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
