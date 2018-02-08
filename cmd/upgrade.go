package cmd

import (
	"github.com/urfave/cli"
	"github.com/odoko-devops/uberstack/uber"
	"github.com/Sirupsen/logrus"
)

func UpgradeCommand() cli.Command {
	return cli.Command{
		Name:   "upgrade",
		Usage:  "Bring all services up",
		Action: uberUpgrade,
		Flags: []cli.Flag{
		},
	}
}

func uberUpgrade(ctx *cli.Context) error {

	logrus.Debug("uberUpgrade")
	uber := uber.Uber{}
	err := uber.Init(ctx)
	if err != nil {
		return err
	}
	uber.Action = "up"
	uber.ActionArguments = []string{"--upgrade"}
	uber.Arguments = ctx.Args()
	err = uber.Execute()
	return err
}
