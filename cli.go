package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/BurntSushi/toml"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if len(flags.Args()) == 0 {
		log.Fatal("argument not found")
		return ExitCodeError
	}

	filePath := flags.Arg(0)
	var data interface{}

	_, err := toml.DecodeFile(filePath, &data)
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	jsonBytes, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	fmt.Println(string(jsonBytes))

	return ExitCodeOK
}
