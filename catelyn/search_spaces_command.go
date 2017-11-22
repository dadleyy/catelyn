package catelyn

import "io"
import "log"
import "net/url"
import "text/tabwriter"
import "github.com/alecthomas/kingpin"

// NewSearchSpacesCommand returns an executable kingpin action.
func NewSearchSpacesCommand(out io.Writer, config *GlobalCLIOptions) Command {
	logger := log.New(out, "", 0)

	cli := searchCLI{
		globals: config,
		Logger:  logger,
		writer:  tabwriter.NewWriter(out, 0, 0, 3, ' ', 10),
	}

	return &cli
}

type searchCLI struct {
	*log.Logger
	globals *GlobalCLIOptions
	query   string
	writer  *tabwriter.Writer
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

	_, paging, e := client.SearchSpaces(c.query)

	if e != nil {
		return e
	}

	c.Printf("has more results: %s", paging.Links.Next)

	return nil
}
