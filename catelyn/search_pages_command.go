package catelyn

import "io"
import "fmt"
import "strings"
import "text/tabwriter"
import "github.com/alecthomas/kingpin"

// NewSearchPagesCommand returns an executable kingpin action for searching pages.
func NewSearchPagesCommand(out io.Writer, config *GlobalCLIOptions) Command {
	cli := &searchPagesCLI{
		output:  out,
		globals: config,
	}
	return cli
}

type searchPagesCLI struct {
	output  io.Writer
	globals *GlobalCLIOptions
	search  ConfluencePageSearchInput
}

func (c *searchPagesCLI) Exec(context *kingpin.ParseContext) error {
	client, e := NewConfluenceClient(c.globals.UserInfo(), c.globals.ConfluenceHost)

	if e != nil {
		return e
	}

	pages, _, e := client.SearchPages(&c.search)

	if e != nil {
		fmt.Fprintf(c.output, "unable to load search results: %v\n", e)
		return nil
	}

	writer := tabwriter.NewWriter(c.output, 10, 2, 3, ' ', 0)

	fmt.Fprintf(c.output, "Pages found on %s (%s):\n", c.globals.ConfluenceHost, c.search.SpaceKey)
	fmt.Fprintf(writer, "Page ID\tPage Title\n")

	for _, p := range pages {
		line := strings.Join([]string{p.ID, p.Title}, "\t")
		fmt.Fprintf(writer, "%s\n", line)
	}

	writer.Flush()

	fmt.Fprintf(c.output, "limit[%d] start[%d]\n", c.search.Limit, c.search.Start)

	return nil
}

func (c *searchPagesCLI) Configure(clause *kingpin.CmdClause) {
	clause.Flag("space", "the key of the space to search for the page in").Short('s').StringVar(&c.search.SpaceKey)
	clause.Flag("title", "the title of the page to search for").Short('t').StringVar(&c.search.Title)
	clause.Flag("limit", "how many results to return").Default("10").Short('l').Uint8Var(&c.search.Limit)
	clause.Flag("start", "how many results to skip").Short('o').Uint8Var(&c.search.Start)
}
