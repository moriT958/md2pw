package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/moriT958/md2pw/internal/converter"
)

type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

func New(outStream, errStream io.Writer) *CLI {
	return &CLI{
		outStream: outStream,
		errStream: errStream,
	}
}

func (c *CLI) Run(args []string) int {
	var outputFile string

	flags := flag.NewFlagSet("md2pw", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.StringVar(&outputFile, "o", "", "output file path (default: stdout)")

	flags.Usage = func() {
		fmt.Fprintf(c.errStream, "Usage: md2pw [options] <file.md>\n\n")
		fmt.Fprintf(c.errStream, "Options:\n")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return 1
	}

	if flags.NArg() < 1 {
		fmt.Fprintln(c.errStream, "Error: input file required")
		flags.Usage()
		return 1
	}

	filename := flags.Arg(0)
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(c.errStream, "Error reading file: %v\n", err)
		return 1
	}

	result, err := converter.Convert(content)
	if err != nil {
		fmt.Fprintf(c.errStream, "Error converting: %v\n", err)
		return 1
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
			fmt.Fprintf(c.errStream, "Error writing to file %s: %v\n", outputFile, err)
			return 1
		}
	} else {
		fmt.Fprint(c.outStream, result)
	}

	return 0
}
