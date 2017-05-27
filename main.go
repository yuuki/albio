package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yuuki/albio/pkg/command"
)

// CLI is the command line object.
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version bool
		status  bool
		attach  bool
		detach  bool
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprint(cli.errStream, helpText)
	}
	flags.BoolVar(&status, "status", false, "")
	flags.BoolVar(&attach, "attach", false, "")
	flags.BoolVar(&detach, "detach", false, "")
	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	if err := flags.Parse(args[1:]); err != nil {
		return 1
	}

	if attach && detach {
		fmt.Fprint(cli.errStream, helpText)
		return 1
	}

	if status {
		if err := command.Status(); err != nil {
			fmt.Fprintln(cli.errStream, err)
			return 2
		}
	} else if attach {
		if err := command.Attach(); err != nil {
			fmt.Fprintln(cli.errStream, err)
			return 2
		}
	} else if detach {
		if err := command.Detach(); err != nil {
			fmt.Fprintln(cli.errStream, err)
			return 2
		}
	} else {
		fmt.Fprint(cli.errStream, helpText)
	}

	if version {
		fmt.Fprintf(cli.errStream, "%s version %s, build %s \n", Name, Version, GitCommit)
		return 0
	}

	return 0
}

var helpText = `
Usage: albio [options]

  A CLI tool to service in/out from AWS Loadbalancer such as ELB/ALB.

Options:
  --status                show loadbalancers information.
  --attach                attach the instance from loadbalancer. It is not possible to specify --detach option if --attach option is specified.
  --detach                detach the instance from loadbalancer. It is not possible to specify --attach option if --detach option is specified.
  --version, -v           print version
`
