package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/apex/log"
	"github.com/fossas/fossa-cli/config"
	"github.com/fossas/fossa-cli/errors"

	"github.com/fossas/fossa-cli/cmd/fossa/display"
	"github.com/fossas/fossa-cli/cmd/fossa/flags"
	"github.com/fossas/fossa-cli/cmd/fossa/setup"
	"github.com/fossas/fossa-cli/cmd/fossa/version"

	"github.com/fossas/fossa-cli/cmd/fossa/cmd/analyze"
	"github.com/fossas/fossa-cli/cmd/fossa/cmd/build"
	initc "github.com/fossas/fossa-cli/cmd/fossa/cmd/init"
	"github.com/fossas/fossa-cli/cmd/fossa/cmd/report"
	"github.com/fossas/fossa-cli/cmd/fossa/cmd/test"
	"github.com/fossas/fossa-cli/cmd/fossa/cmd/update"
	"github.com/fossas/fossa-cli/cmd/fossa/cmd/upload"
)

// App describes the main application.
var App = cli.App{
	Name:                 "fossa-cli",
	Usage:                "Fast, portable and reliable dependency analysis (https://github.com/fossas/fossa-cli/)",
	Version:              version.String(),
	Action:               Run,
	EnableBashCompletion: true,
	Flags: flags.Combine(
		initc.Cmd.Flags,
		analyze.Cmd.Flags,
		flags.WithGlobalFlags(nil),
	),
	Commands: []cli.Command{
		initc.Cmd,
		build.Cmd,
		analyze.Cmd,
		upload.Cmd,
		report.Cmd,
		test.Cmd,
		update.Cmd,
	},
}

func main() {
	err := App.Run(os.Args)
	if err != nil {
		switch e := err.(type) {
		case *errors.Error:
			display.Error(errors.Render(e.Error(), errors.MessageArgs{
				Invocation: strings.Join(os.Args, " "),
				LogFile:    display.File(),
			}))
			os.Exit(e.ExitCode)
		default:
			log.Errorf("Error: %#v", err.Error())
			os.Exit(1)
		}
	}
}

func Run(ctx *cli.Context) error {
	err := setup.SetContext(ctx)
	if err != nil {
		return err
	}

	if config.APIKey() == "" && !ctx.Bool(analyze.ShowOutput) {
		fmt.Printf("Incorrect Usage. FOSSA_API_KEY must be set as an environment variable or provided in .fossa.yml\n\n")
		log.Fatalf("No API KEY provided")
	}

	err = initc.Run(ctx)
	if err != nil {
		return err
	}

	err = analyze.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}
