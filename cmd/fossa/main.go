package main

import (
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/urfave/cli"

	"github.com/fossas/fossa-cli/api/fossa"
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
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	err := App.Run(os.Args)
	if err != nil {
		log.WithError(err).Error("fatal error")
		switch e := err.(type) {
		case *errors.Error:
			display.Error(errors.Render(e.Error(), errors.MessageArgs{
				Invocation: strings.Join(os.Args, " "),
				Endpoint:   config.Endpoint(),
				LogFile:    display.File(),
			}))
			if e.ExitCode != 0 {
				os.Exit(e.ExitCode)
			} else {
				os.Exit(1)
			}
		default:
			display.Error(Error(err))
			os.Exit(1)
		}
	}
}

func Run(ctx *cli.Context) error {
	err := setup.SetContext(ctx)
	if err != nil {
		return err
	}

	if !ctx.Bool(analyze.ShowOutput) {
		err = fossa.SetAPIKey(config.APIKey())
		if err != nil {
			return err
		}
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
