package cmd

import (
	"github.com/urfave/cli"
	"github.com/odoko-devops/uberstack/uber"
	"github.com/Sirupsen/logrus"
)

func ConfirmCommand() cli.Command {
	return cli.Command{
		Name:   "confirm",
		Usage:  "Confirm previous upgrade",
		Action: uberConfirm,
		Flags: []cli.Flag{
		},
	}
}

func uberConfirm(ctx *cli.Context) error {

	logrus.Debug("uberConfirm")
	uber := uber.Uber{}
	err := uber.Init(ctx)
	if err != nil {
		return err
	}
	uber.Action = "up"
	uber.ActionArguments = []string{"--confirm-upgrade"}
	uber.Arguments = ctx.Args()
	err = uber.Execute()
	return err
}
