//go:generate go-bindata -pkg main -o credits.go CREDITS
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yuuki/albio/pkg/command"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 10 + iota
)

var (
	creditsText = string(MustAsset("CREDITS"))
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
		return exitCodeErr
	}

	var err error

	switch args[1] {
	case "status":
		err = cli.doStatus(args[2:])
	case "attach":
		err = cli.doAttach(args[2:])
	case "detach":
		err = cli.doDetach(args[2:])
	case "--version":
		fmt.Fprintf(cli.errStream, "%s version %s, build %s, date %s \n", name, version, commit, date)
		return exitCodeOK
	case "-h", "--help":
		fmt.Fprint(cli.errStream, helpText)
		return exitCodeOK
	case "--credits":
		fmt.Fprint(cli.errStream, creditsText)
		return exitCodeOK
	default:
		fmt.Fprint(cli.errStream, helpText)
		return exitCodeErr
	}

	if err != nil {
		if os.Getenv("ALBIO_DEBUG") != "" {
			fmt.Fprintf(cli.errStream, "%+v\n", err)
		} else {
			fmt.Fprintln(cli.errStream, err)
		}
		return exitCodeErr
	}

	return 0
}

var helpText = `Usage: albio [options]

  A CLI tool to service in/out from AWS Loadbalancer such as ALB/NLB.

Commands:
  status	show loadbalancers information.
  attach    attach the EC2 instance to the loadbalancers.
  detach 	detach the EC2 instance from the loadbalancers.

Options:
  --version, -v		print version
  --credits         print credits
  --help, -h		print help
`

func (cli *CLI) prepareFlags(help string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
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

var attachHelpText = `
Usage: albio attach [options]

attach the EC2 instance to the loadbalancers.

Options:
  --instance-id, -i		specify EC2 instance id
  --self				use localhost as EC2 instance
  --loadbalancer, -l	specify ALB/NLB name or ARN
`

func (cli *CLI) doAttach(args []string) error {
	var param command.AttachParam
	flags := cli.prepareFlags(attachHelpText)
	flags.StringVar(&param.InstanceID, "i", "", "")
	flags.StringVar(&param.InstanceID, "instance-id", "", "")
	flags.BoolVar(&param.Self, "self", false, "")
	flags.StringVar(&param.LoadBalancerName, "l", "", "")
	flags.StringVar(&param.LoadBalancerName, "loadbalancer", "", "")
	if err := flags.Parse(args); err != nil {
		return err
	}
	return command.Attach(&param)
}

var detachHelpText = `
Usage: albio detach [options]

detach the EC2 instance from the loadbalancers.

Options:
  --instance-id, -i		specify EC2 instance id
  --self            	use localhost as EC2 instance
`

func (cli *CLI) doDetach(args []string) error {
	var param command.DetachParam
	flags := cli.prepareFlags(detachHelpText)
	flags.StringVar(&param.InstanceID, "i", "", "")
	flags.StringVar(&param.InstanceID, "instance-id", "", "")
	flags.BoolVar(&param.Self, "self", false, "")
	if err := flags.Parse(args); err != nil {
		return err
	}
	return command.Detach(&param)
}
