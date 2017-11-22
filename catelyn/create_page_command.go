package catelyn

import "io"
import "fmt"
import "io/ioutil"
import "github.com/chzyer/readline"
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
		title    string
		space    string
		filename string
	}
}

func (c *createPageCLI) Configure(clause *kingpin.CmdClause) {
	clause.Flag("space", "the space to create the new page under").Short('s').StringVar(&c.page.space)
	clause.Flag("title", "the title of the page to create").Short('t').StringVar(&c.page.title)
	clause.Flag("content", "the file that should be read for the content").Short('f').StringVar(&c.page.filename)
	clause.Action(c.prompt)
}

func (c *createPageCLI) Exec(context *kingpin.ParseContext) error {
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
		if _, e := fmt.Fscanln(c.reader, &c.page.title); e != nil {
			return e
		}
	}

	if c.page.filename != "" {
		return nil
	}

	completer := readline.NewPrefixCompleter(
		readline.PcItemDynamic(c.listFiles("./")),
	)

	rl, e := readline.NewEx(&readline.Config{
		Prompt:       "filename: ",
		AutoComplete: completer,
	})

	if e != nil {
		return e
	}

	defer rl.Close()

	c.page.filename, e = rl.Readline()

	if e != nil {
		return e
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
