package dockerfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Statement interface{}

type Dockerfile struct {
	Source     string
	Statements []Statement
	From       *From
}

type Comment struct {
	Line  int
	Lines string
}

type From struct {
	Line  int
	Image string
}

type Maintainer struct {
	Line int
	Name string
}

type Run struct {
	Line    int
	Command string
}

type Label struct {
	Line  int
	Pairs string
}

type Expose struct {
	Line  int
	Ports string
}

type Env struct {
	Line  int
	Pairs string
}

type Add struct {
	Line  int
	Paths string
}

type Copy struct {
	Line  int
	Paths string
}

type Entrypoint struct {
	Line    int
	CmdLine string
}

type Volume struct {
	Line  int
	Names string
}

type User struct {
	Line int
	Name string
}

type Workdir struct {
	Line int
	Path string
}

type Arg struct {
	Line int
	Name string
}

type Onbuild struct {
	Line  			int
	Instruction string
}

type Stopsignal struct {
	Line   int
	Signal string
}

type Healthcheck struct {
	Line    int
	Command string
}

type Shell struct {
	Line    int
	CmdLine string
}

type Cmd struct {
	Line    int
	CmdLine string
}

var (
	statements = map[string]func(file *Dockerfile, line int, tokens *Tokens) Statement{
		"FROM": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			file.From = &From{Line: line, Image: tokens.NextLine()}
			return file.From
		},
		"MAINTAINER": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Maintainer{Line: line, Name: tokens.NextLine()}
		},
		"RUN": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Run{Line: line, Command: tokens.NextLine()}
		},
		"CMD": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Cmd{Line: line, CmdLine: tokens.NextLine()}
		},
		"LABEL": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Label{Line: line, Pairs: tokens.NextLine()}
		},
		"EXPOSE": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Expose{Line: line, Ports: tokens.NextLine()}
		},
		"ENV": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Env{Line: line, Pairs: tokens.NextLine()}
		},
		"ADD": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Add{Line: line, Paths: tokens.NextLine()}
		},
		"COPY": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Copy{Line: line, Paths: tokens.NextLine()}
		},
		"ENTRYPOINT": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Entrypoint{Line: line, CmdLine: tokens.NextLine()}
		},
		"VOLUME": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Volume{Line: line, Names: tokens.NextLine()}
		},
		"USER": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &User{Line: line, Name: tokens.NextLine()}
		},
		"WORKDIR": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Workdir{Line: line, Path: tokens.NextLine()}
		},
		"ARG": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Arg{Line: line, Name: tokens.NextLine()}
		},
		"ONBUILD": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Onbuild{Line: line, Instruction: tokens.NextLine()}
		},
		"STOPSIGNAL": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Stopsignal{Line: line, Signal: tokens.NextLine()}
		},
		"HEALTHCHECK": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Healthcheck{Line: line, Command: tokens.NextLine()}
		},
		"SHELL": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Shell{Line: line, CmdLine: tokens.NextLine()}
		},
		"#": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Comment{Line: line, Lines: tokens.NextComment()}
		},
	}
)

// Parses a dockerfile from a reader. Returns an error if
// an unknown token is encountered.
//
// See https://docs.docker.com/engine/reference/builder/
func Parse(input io.Reader, file *Dockerfile, source ...string) (err error) {
	tokens := NewTokens(input)

	for tokens.HasNext {
		token := tokens.NextToken()

		if "" == token {
			continue
		} else if statement, ok := statements[token]; ok {
			file.Statements = append(file.Statements, statement(file, tokens.Line, tokens))
		} else {
			return fmt.Errorf("Cannot handle token `%s` on line %d", token, tokens.Line)
		}
	}

	if len(source) > 0 {
		file.Source = source[0]
	} else {
		file.Source = fmt.Sprintf("%T", input)
	}

	return nil
}

// Parses a dockerfile from a file. Returns an error if
// the file cannot be opened, is a directory or when parsing
// encounters an error
func ParseFile(name string, file *Dockerfile) (err error) {
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("The given file `%s` is a directory\n", name)
	}

	input, err := os.Open(name)
	if err != nil {
		return err
	}

	defer input.Close()
	return Parse(bufio.NewReader(input), file, name)
}

// Extends parser. Example:
//
// 		type Include struct {
// 		  Line      int
// 		  Reference string
// 		}
//
// 		Extend("INCLUDE", func(file *Dockerfile, line int, tokens *Tokens) Statement {
// 		  return &Include{Line: line, Reference: tokens.NextLine()}
// 		})
//
func Extend(name string, extension func(file *Dockerfile, line int, tokens *Tokens) Statement) {
	statements[name] = extension
}