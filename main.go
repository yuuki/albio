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
	if len(args) <= 1 {
		fmt.Fprint(cli.errStream, helpText)
		return 2
	}

	var err error

	switch args[1] {
	case "status":
		err = cli.doStatus(args[2:])
	case "attach":
		err = cli.doAttach(args[2:])
	case "detach":
		err = cli.doDetach(args[2:])
	case "-v", "--version":
		fmt.Fprintf(cli.errStream, "%s version %s, build %s \n", Name, Version, GitCommit)
		return 0
	case "-h", "--help":
		fmt.Fprint(cli.errStream, helpText)
	default:
		fmt.Fprint(cli.errStream, helpText)
	}

	if err != nil {
		fmt.Fprintln(cli.errStream, err)
		return 2
	}

	return 0
}

var helpText = `
Usage: albio [options]

  A CLI tool to service in/out from AWS Loadbalancer such as ELB/ALB.

Commands:
  status	show loadbalancers information.
  attach        attach the instance from loadbalancer. It is not possible to specify --detach option if --attach option is specified.
  detach        detach the instance from loadbalancer. It is not possible to specify --attach option if --detach option is specified.

Options:
  --version, -v		print version
  --help, -h            print help
`

func (cli *CLI) prepareFlags(help string) *flag.FlagSet {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprint(cli.errStream, help)
	}
	return flags
}

var statusHelpText = `
Usage: albio status [options]

show loadbalancers information.

Options:
  --instance-id, -i	specify EC2 instance id
`

func (cli *CLI) doStatus(args []string) error {
	var param command.StatusParam
	flags := cli.prepareFlags(statusHelpText)
	flags.StringVar(&param.InstanceID, "i", "", "")
	flags.StringVar(&param.InstanceID, "instance-id", "", "")
	if err := flags.Parse(args); err != nil {
		return err
	}
	return command.Status(&param)
}

func (cli *CLI) doAttach(args []string) error {
	flags := cli.prepareFlags(helpText)
	if err := flags.Parse(args); err != nil {
		return err
	}
	return command.Attach()
}

func (cli *CLI) doDetach(args []string) error {
	flags := cli.prepareFlags(helpText)
	if err := flags.Parse(args); err != nil {
		return err
	}
	return command.Detach()
}
