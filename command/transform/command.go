package transform

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/command/clean"
	"github.com/tueftler/doget/dockerfile"
)

type TransformCommand struct {
	command.Command
	flags *flag.FlagSet
}

// Creates new transform command instance
func NewCommand(name string) *TransformCommand {
	return &TransformCommand{flags: flag.NewFlagSet(name, flag.ExitOnError)}
}

func open(output string) (io.Writer, error) {
	if output == "-" {
		return os.Stdout, nil
	} else {
		return os.Create(output)
	}
}

// Runs transform command
func (c *TransformCommand) Run(parser *dockerfile.Parser, args []string) error {
	input := c.flags.String("in", "Dockerfile.in", "Input. Use - for standard input")
	output := c.flags.String("out", "Dockerfile", "Output. Use - for standard output")
	performClean := c.flags.Bool("clean", false, "Remove vendor directory after transformation")
	noCache := c.flags.Bool("no-cache", false, "Do not use cache")
	c.flags.Parse(args)

	fmt.Fprintf(os.Stderr, "> Running transform(%q -> %q)\n", *input, *output)

	if *performClean {
		defer clean.NewCommand("clean").Run(parser, args)
	}

	// Open output
	out, err := open(*output)
	if err != nil {
		return err
	}

	// Transform
	var buf bytes.Buffer
	transformation := Transformation{Input: *input, Output: &buf, UseCache: !*noCache}
	if err := transformation.Run(parser); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Done\n\n")

	// Result
	fmt.Fprintf(out, buf.String())
	return nil
}
