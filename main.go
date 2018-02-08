package main

import (
	"os"

	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/odoko-devops/uberstack/cmd"
	"github.com/urfave/cli"
)

var VERSION = "dev"

var AppHelpTemplate = `{{.Usage}}

Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

Version: {{.Version}}
{{if .Flags}}
Options:
  {{range .Flags}}{{if .Hidden}}{{else}}{{.}}
  {{end}}{{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

var CommandHelpTemplate = `{{.Usage}}
{{if .Description}}{{.Description}}{{end}}
Usage: uber [global options] {{.Name}} {{if .Flags}}[OPTIONS] {{end}}{{if ne "None" .ArgsUsage}}{{if ne "" .ArgsUsage}}{{.ArgsUsage}}{{else}}[arg...]{{end}}{{end}}

{{if .Flags}}Options:{{range .Flags}}
	 {{.}}{{end}}{{end}}
`

func main() {
	if err := mainErr(); err != nil {
		logrus.Fatal(err)
	}
}

func mainErr() error {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate

	app := cli.NewApp()
	app.Name = "uber"
	app.Usage = "UberStack CLI, managing applications."
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Version = VERSION
	app.Author = "Odoko Ltd."
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug logging",
		},
		cli.StringFlag{
			Name:   "config,c",
			Value: "./uber.yml",
			Usage:  "Client configuration file (default ./uber.yml)",
			EnvVar: "UBER_CONFIG",
		},
	}
	app.Commands = []cli.Command{
		cmd.UpgradeCommand(),
	}
	parsed, err := parseArgs(os.Args)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	return app.Run(parsed)
}

var singleAlphaLetterRegxp = regexp.MustCompile("[a-zA-Z]")

func parseArgs(args []string) ([]string, error) {
	result := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 1 {
			for i, c := range arg[1:] {
				if string(c) == "=" {
					if i < 1 {
						return nil, errors.New("invalid input with '-' and '=' flag")
					}
					result[len(result)-1] = result[len(result)-1] + arg[i+1:]
					break
				} else if singleAlphaLetterRegxp.MatchString(string(c)) {
					result = append(result, "-"+string(c))
				} else {
					return nil, errors.Errorf("invalid input %v in flag", string(c))
				}
			}
		} else {
			result = append(result, arg)
		}
	}
	return result, nil
}
