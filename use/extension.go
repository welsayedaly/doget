package use

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/tueftler/doget/dockerfile"
)

// Statement represents a single USE statement
type Statement struct {
	Context   *Context
	Line      int
	Reference string
}

// Context represents the use context
type Context struct {
	Repositories map[string]map[string]string
}

// Origin represents the parsed components of a USE reference
type Origin struct {
	Host    string
	Vendor  string
	Name    string
	Version string
	Dir     string
	Uri     string
}

// New creates a USE instruction backed by the given repositories
func New(repositories map[string]map[string]string) *Context {
	return &Context{Repositories: repositories}
}

// String creates a string representation of an origin
func (o *Origin) String() string {
	str := o.Host + "/" + o.Vendor + "/" + o.Name
	if "" != o.Dir {
		str += "/" + o.Dir
	}
	if "" != o.Version {
		str += ":" + o.Version
	}
	return str
}

// Emit writes the USE statement
func (s *Statement) Emit(out io.Writer) {
	dockerfile.EmitInstruction(out, "USE", s.Reference)
}

// Origin parses origin from reference
func (s *Statement) Origin() (origin *Origin, err error) {
	var parsed []string

	origin = &Origin{}

	// Version
	pos := strings.LastIndex(s.Reference, ":")
	if pos == -1 {
		parsed = strings.Split(s.Reference, "/")
		origin.Version = "master"
	} else {
		parsed = strings.Split(s.Reference[0:pos], "/")
		origin.Version = s.Reference[pos+1 : len(s.Reference)]
	}

	origin.Host = parsed[0]
	origin.Vendor = parsed[1]
	origin.Name = parsed[2]

	// Subdirectory
	if len(parsed) > 3 {
		origin.Dir = strings.Join(parsed[3:len(parsed)], "/")
	} else {
		origin.Dir = ""
	}

	// Compile URL
	if repository, ok := s.Context.Repositories[origin.Host]; ok {
		template, err := template.New(origin.Host).Parse(repository["url"])
		if err != nil {
			return nil, err
		}

		var uri bytes.Buffer
		if err := template.Execute(&uri, *origin); err != nil {
			return nil, err
		}

		origin.Uri = uri.String()
		return origin, nil
	} else {
		return nil, fmt.Errorf("No repository %s", origin.Host)
	}
}

// Extension func for parser
func (c *Context) Extension(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
	return &Statement{Context: c, Line: line, Reference: strings.TrimSpace(tokens.NextLine())}
}
