package transform

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/config"
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

// Runs transform command
func (c *TransformCommand) Run(parser *dockerfile.Parser, args []string) error {
	input := c.flags.String("in", "Dockerfile.in", "Input. Use - for standard input")
	output := c.flags.String("out", "Dockerfile", "Output. Use - for standard output")
	performClean := c.flags.Bool("clean", false, "Remove "+config.Vendordir+" directory after transformation")
	noCache := c.flags.Bool("no-cache", false, "Do not use cache")
	c.flags.Parse(args)

	if *performClean {
		defer os.RemoveAll(config.Vendordir)
	}

	storage := config.Vendordir + ".zip"
	if _, err := os.Stat(storage); err == nil {
		fmt.Fprint(os.Stderr, "Preparing...")
		if err := unzip(storage, ".", strings.NewReplacer()); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, " OK")
	}

	// Transform
	var buf bytes.Buffer
	transformation := Transformation{Input: *input, Output: &buf, UseCache: !*noCache}
	err := transformation.Run(parser)

	if err == nil {
		fmt.Fprint(os.Stderr, "Caching...")
		if err := mkzip(config.Vendordir, storage); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, " Done")
	}

	if err != nil {
		return err
	}

	// Result
	if *output == "-" {
		fmt.Println(buf.String())
	} else {
		out, err := os.Create(*output)
		if err != nil {
			return err
		}

		defer out.Close()
		fmt.Fprintf(out, buf.String())
	}

	return nil
}
