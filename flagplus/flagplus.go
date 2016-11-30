// Copyright 2011 The Go Authors
// Copyright 2013 Jonas mg
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

// Based in code from 'http://code.google.com/p/go/source/browse/src/cmd/go/main.go'

// Package flagplus provides management of sub-commands to command-line programs.
// It is based in code from "go tool".
package flagplus

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// A Subcommand is an implementation of a sub-command.
// Every command gets an extra flag, help, which shows its documentation.
type Subcommand struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Subcommand, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the "<program> help" output.
	Short string

	// Long is the long message shown in the "<program> help <this-command>" output.
	Long string

	// Flag is a set of flags specific for this command; it is set using method
	// (*Subcommand).AddFlags
	FlagSet flag.FlagSet

	// CustomFlags indicates that the command will do its own flag parsing.
	CustomFlags bool
}

// AddFlags looks up the flags in the global flag.FlagSet and they are added
// to the sub-command, if they are found.
func (s *Subcommand) AddFlags(names ...string) *Subcommand {
	invalidNames := make([]string, 0)
	hasError := false

	for _, v := range names {
		_flag := flag.Lookup(v)
		if _flag == nil {
			invalidNames = append(invalidNames, v)
			hasError = true
		} else if !hasError {
			s.FlagSet.Var(_flag.Value, _flag.Name, _flag.Usage)
		}
	}

	if hasError {
		errMsg := "flag does not exist"
		if len(invalidNames) != 1 {
			errMsg = "flags do not exist"
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", errMsg, strings.Join(invalidNames, ", "))
		os.Exit(2)
	}
	return s
}

// Name returns the command's name: the first word in the usage line.
func (s *Subcommand) Name() string {
	name := s.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (s *Subcommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s %s\n\n", os.Args[0], s.UsageLine)
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (s *Subcommand) Runnable() bool {
	return s.Run != nil
}

// * * *

// Command represents the structure of a main command with sub-commands.
// Every command gets an extra flag, <help documentation>, which enables to
// generate documentation of all sub-commands in 'doc.go'.
type Command struct {
	Description   string // The program's description shown in the "<program> help" output.
	Subcommands   []*Subcommand
	globalFlags   []string
	ErrorHandling flag.ErrorHandling
}

// NewCommand creates a new command with a default ErrorHandling to
// "flag.ExitOnError".
func NewCommand(description string, cmds ...*Subcommand) *Command {
	return &Command{
		Description:   description,
		Subcommands:   cmds,
		ErrorHandling: flag.ExitOnError,
	}
}

// AddGlobalFlags sets the global flags which will be parsed using
// "*Command.Parse()".
func (c *Command) AddGlobalFlags(names ...string) *Command {
	invalidNames := make([]string, 0)
	hasError := false

	for _, v := range names {
		if flag.Lookup(v) == nil {
			invalidNames = append(invalidNames, v)
			hasError = true
		}
	}

	if hasError {
		errMsg := "flag does not exist"
		if len(invalidNames) != 1 {
			errMsg = "flags do not exist"
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", errMsg, strings.Join(invalidNames, ", "))
		os.Exit(2)
	}

	c.globalFlags = names
	return c
}

// HasGlobalFlags tests whether the command has global flags.
func (c *Command) HasGlobalFlags() bool {
	if len(c.globalFlags) != 0 {
		return true
	}
	return false
}

// GetGlobalFlags returns the global flags.
func (c *Command) GetGlobalFlags() []string { return c.globalFlags }

// Parse parses both command and flag definitions from the argument list.
// Also, the global flags are added to each sub-command, if any.
func (c *Command) Parse() {
	var err error

	flag.Usage = c.Usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		err = fmt.Errorf("%s\n", c.Description)
	} else if args[0] == "help" {
		if err = c.help(args[1:]); err == nil {
			return
		}
	} else {
		hasGlobalFlags := false
		if len(c.globalFlags) != 0 {
			hasGlobalFlags = true
		}

		for _, subc := range c.Subcommands {
			if subc.Name() == args[0] && subc.Run != nil {
				subc.FlagSet.Usage = func() { subc.Usage() }

				if hasGlobalFlags {
					for _, v := range c.globalFlags {
						_flag := flag.Lookup(v)
						subc.FlagSet.Var(_flag.Value, _flag.Name, _flag.Usage)
					}
				}

				if subc.CustomFlags {
					args = args[1:]
				} else {
					subc.FlagSet.Parse(args[1:])
					args = subc.FlagSet.Args()
				}

				subc.Run(subc, args)
				return
			}
		}
		err = fmt.Errorf("Unknown subcommand %q.  Run `%s help` for usage.\n",
			args[0], os.Args[0])
	}

	switch c.ErrorHandling {
	case flag.ContinueOnError:
		fmt.Fprint(os.Stderr, err)
	case flag.ExitOnError:
		fmt.Fprint(os.Stderr, err)
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	}
}

func (c *Command) Usage() {
	c.printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the "help" command.
// "help documentation" generates documentation of all commands in 'doc.go'.
func (c *Command) help(args []string) error {
	if len(args) == 0 { // Succeeded at "<program> help".
		c.printUsage(os.Stdout)
		return nil
	}
	if len(args) != 1 { // Failed at "<program> help".
		return fmt.Errorf("Usage: %s help command\n\nToo many arguments given.\n", os.Args[0])
	}

	arg := args[0]

	// "<program> help documentation" generates 'doc.go'.
	if arg == "documentation" {
		buf := new(bytes.Buffer)
		c.printUsage(buf)
		usage := &Subcommand{Long: buf.String()}

		file, err := os.OpenFile("doc.go", os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
		tmpl(file, documentationTemplate, append([]*Subcommand{usage}, c.Subcommands...))
		return file.Close()
	}

	for _, subc := range c.Subcommands {
		if subc.Name() == arg { // Succeeded at "<program> help <cmd>".
			tmpl(os.Stdout, helpTemplate, subc)
			return nil
		}
	}
	// Failed at "<program> help <cmd>"
	return fmt.Errorf("Unknown help topic %q.  Run `%s help` for usage.\n", arg, os.Args[0])
}

func (c *Command) printUsage(w io.Writer) { tmpl(w, usageTemplate, c) }

// == Templates
//

var usageTemplate = `{{.Description}}

Usage:
      {{program}}{{if .HasGlobalFlags}} [global flags]{{end}} command [flags] [arguments]

## Commands
{{range .Subcommands}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "{{program}} help [command]" for more information about a command.
{{if hasExtraTopic .Subcommands}}
Additional help topics:
{{range .Subcommands}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "{{program}} help [topic]" for more information about that topic.
{{end}}
{{if .HasGlobalFlags}}## Global flags
{{printGlobFlags .GetGlobalFlags}}{{end}}`

var helpTemplate = `{{if .Runnable}}Usage: {{program}} {{.UsageLine}}

{{end}}{{.Long | trim}}
{{if hasFlags .FlagSet}}
Flags:
{{printDefaults .FlagSet}}{{end}}
`

var documentationTemplate = `// DO NOT EDIT THIS FILE. GENERATED BY "{{cmdLine}}".
// Edit the documentation in other files and re-run it to generate this one.

/*
{{range .}}{{if .Short}}***

{{.Short | capitalize}}

{{end}}{{if .Runnable}}Usage: {{program}} {{.UsageLine}}

{{end}}{{.Long | trim}}

{{end}}*/
package main
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"capitalize":    capitalize,
		"hasExtraTopic": hasExtraTopic,
		"hasFlags":      hasFlags,
		"trim":          strings.TrimSpace,

		"cmdLine": func() string { return strings.Join(os.Args, " ") },
		"program": func() string { return os.Args[0] },

		"printDefaults": func(f flag.FlagSet) string {
			f.SetOutput(w)
			f.PrintDefaults()
			return ""
		},
		"printGlobFlags": func(names []string) string {
			flag.VisitAll(func(f *flag.Flag) {
				found := false
				for _, v := range names {
					if v == f.Name {
						found = true
						break
					}
				}
				if !found {
					return
				}

				format := ""

				valueType := reflect.TypeOf(f.Value)
				baseValueType := valueType.Elem() // we follow the pointer

				switch baseValueType.Kind() {
				case reflect.String:
					format = "\n  -%s=%q: %s" // put quotes on the value
				default:
					format = "\n  -%s=%s: %s"
				}

				fmt.Fprintf(w, format, f.Name, f.DefValue, f.Usage)
			})

			fmt.Fprint(w, "\n\n")
			return ""
		},
	})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func hasExtraTopic(cmds []*Subcommand) bool {
	for _, v := range cmds {
		if !v.Runnable() {
			return true
		}
	}
	return false
}

func hasFlags(f flag.FlagSet) bool {
	nFlags := 0

	f.VisitAll(func(*flag.Flag) {
		nFlags++
	})
	if nFlags != 0 {
		return true
	}
	return false
}
