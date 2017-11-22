package catelyn

import "io"
import "fmt"
import "bufio"
import "strings"
import "io/ioutil"
import "github.com/alecthomas/kingpin"

// NewCreatePageCommand returns an implementation of the command interface used for creating confluence pages.
func NewCreatePageCommand(output io.Writer, input io.Reader, config *GlobalCLIOptions) Command {
	cli := &createPageCLI{
		output:  output,
		reader:  input,
		globals: config,
	}

	return cli
}

type createPageCLI struct {
	output  io.Writer
	reader  io.Reader
	globals *GlobalCLIOptions
	page    struct {
		title  string
		space  string
		parent string
	}
}

func (c *createPageCLI) Configure(clause *kingpin.CmdClause) {
	clause.Flag("space", "the space to create the new page under").Short('s').StringVar(&c.page.space)
	clause.Flag("title", "the title of the page to create").Short('t').StringVar(&c.page.title)
	clause.Flag("parent", "the id of the page to nest the created page under").Short('a').StringVar(&c.page.parent)
	clause.Action(c.prompt)
}

func (c *createPageCLI) Exec(context *kingpin.ParseContext) error {
	fmt.Fprintf(c.output, "creating \"%s\" in %s\n", c.page.title, c.page.space)

	client, e := NewConfluenceClient(c.globals.UserInfo(), c.globals.ConfluenceHost)

	if e != nil {
		return e
	}

	input := ConfluencePageCreationInput{
		Title:    c.page.title,
		SpaceKey: c.page.space,
		ParentID: c.page.parent,
	}

	if e := client.CreatePage(input); e != nil {
		fmt.Fprintf(c.output, "unable to create page: %v\n", e)
		return nil
	}

	fmt.Fprintf(c.output, "successfully created page!\n")
	return nil
}

func (c *createPageCLI) prompt(context *kingpin.ParseContext) error {
	if c.page.space == "" {
		fmt.Fprintf(c.output, "space key: ")

		if _, e := fmt.Fscanln(c.reader, &c.page.space); e != nil {
			return e
		}
	}

	if c.page.title == "" {
		fmt.Fprintf(c.output, "title: ")
		reader := bufio.NewReader(c.reader)
		token, e := reader.ReadString('\n')

		if e != nil {
			return e
		}

		c.page.title = strings.TrimSuffix(token, "\n")
	}

	return nil
}

func (c *createPageCLI) listFiles(dir string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(dir)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}
