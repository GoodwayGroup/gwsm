package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gw-kube-env/cmd"
	"gw-kube-env/info"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	app := &cli.App{
		Name:     info.AppName,
		Version:  info.AppVersion,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Derek Smith",
				Email: "dsmith@goodwaygroup.com",
			},
		},
		Copyright:            "(c) 2020 Goodway Group",
		HelpName:             info.AppName,
		Usage:                "interact with config map and secret manager variables",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "region",
				Usage:       "AWS region",
				Required:    false,
				Value:       "us-east-1",
				DefaultText: "us-east-1",
				EnvVars:     []string{"AWS_REGION"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version info",
				Action: func(c *cli.Context) error {
					fmt.Printf("%s %s (%s/%s)\n", info.AppName, info.AppVersion, runtime.GOOS, runtime.GOARCH)
					return nil
				},
			},
			{
				Name:      "view",
				Usage:     "View Env local or within a namespace",
				UsageText: info.ViewCommandHelpText,
				Subcommands: []*cli.Command{
					{
						Name:      "local",
						Usage:     "View values based on local settings",
						UsageText: "View the current environment variables based on the LOCAL configurations for a given configmap and secrets.yml",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secrets-path",
								Aliases:  []string{"s"},
								Usage:    "Path to secrets.yml",
								Required: false,
								Value:    ".docker/secrets.yml",
							},
							&cli.StringFlag{
								Name:     "configmap-path",
								Aliases:  []string{"m"},
								Usage:    "Path to configmap.yaml",
								Required: true,
							},
						},
						Action: cmd.ViewLocalEnv,
					},
					{
						Name:      "kube",
						Usage:     "View values configured withing a namesapce",
						UsageText: "View the current environment variables of a given command running on a pod within a namespace.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "namespace",
								Aliases:  []string{"n"},
								Usage:    "Kube Namespace",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "cmd",
								Usage:    "Command to grep for",
								Required: false,
								Value:    "node",
							},
							&cli.StringFlag{
								Name:     "filter",
								Usage:    "Prefix used to filter unwanted env vars from printing",
								Required: false,
								Value:    "npm_",
							},
						},
						Action: cmd.ViewNamespaceEnv,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
