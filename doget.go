package main

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/command/dump"
	"github.com/tueftler/doget/command/transform"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
	"github.com/tueftler/doget/use"
	"os"
)

var (
	parser   = dockerfile.NewParser().Extend("USE", use.Extension)
	commands = map[string]command.Command{
		"dump":      dump.NewCommand("dump"),
		"transform": transform.NewCommand("transform"),
	}
)

func configuration(file string) (*config.Configuration, error) {
	if file == "" {
		return config.Default().Merge(config.SearchPath()...)
	} else {
		return config.Empty().MustMerge(file)
	}
}

func main() {
	var (
		cmdName    = flag.String("#1", "", "Command, one of [dump, transform]")
		configFile = flag.String("config", "", "Configuration file to use")
	)
	flag.Parse()

	configuration, err := configuration(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	*cmdName = flag.Arg(0)
	if delegate, ok := commands[*cmdName]; ok {
		args := flag.Args()
		if err := delegate.Run(configuration, parser, args[1:len(args)]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot handle command %q\n", *cmdName)
		flag.Usage()
		os.Exit(2)
	}
}
