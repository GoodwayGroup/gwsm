package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gwsm/cmd"
	"gwsm/info"
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
						Aliases:   []string{"l"},
						Usage:     "View values based on local settings",
						UsageText: info.ViewLocalCommandHelpText,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secrets",
								Aliases:  []string{"s"},
								Usage:    "Path to secrets.yml",
								Required: false,
								Value:    ".docker/secrets.yml",
							},
							&cli.StringFlag{
								Name:     "configmap",
								Aliases:  []string{"c"},
								Usage:    "Path to configmap.yaml",
								Required: true,
							},
						},
						Action: cmd.ViewLocalEnv,
					},
					{
						Name:      "namespace",
						Aliases:   []string{"ns"},
						Usage:     "View values configured withing a namespace",
						UsageText: info.ViewNamespaceCommandHelpText,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "namespace",
								Aliases:  []string{"n"},
								Usage:    "Kube Namespace list Pods from",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "cmd",
								Usage:    "Command to inspect",
								Required: false,
								Value:    "node",
							},
							&cli.StringFlag{
								Name:     "filter-prefix",
								Aliases:  []string{"f"},
								Usage:    "List of prefixes (csv) used to filter values from display. Set to \"\" to remove any filters.",
								Required: false,
								Value:    "npm_,KUBERNETES_,API_PORT",
							},
							&cli.StringFlag{
								Name:     "exclude",
								Usage:    "List (csv) of specific env vars to exclude values from display. Set to \"\" to remove any exclusions.",
								Required: false,
								Value:    "PATH,SHLVL,HOSTNAME",
							},
						},
						Action: cmd.ViewNamespaceEnv,
					},
					{
						Name:      "diff",
						Aliases:   []string{"d"},
						Usage:     "View diff of local vs. namespace",
						UsageText: info.ViewDiffCommandHelpText,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secrets",
								Aliases:  []string{"s"},
								Usage:    "Path to secrets.yml",
								Required: false,
								Value:    ".docker/secrets.yml",
							},
							&cli.StringFlag{
								Name:     "configmap",
								Aliases:  []string{"c"},
								Usage:    "Path to configmap.yaml",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "namespace",
								Aliases:  []string{"n"},
								Usage:    "Kube Namespace list Pods from",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "cmd",
								Usage:    "Command to inspect",
								Required: false,
								Value:    "node",
							},
							&cli.StringFlag{
								Name:     "filter-prefix",
								Aliases:  []string{"f"},
								Usage:    "List of prefixes (csv) used to filter values from display. Set to \"\" to remove any filters.",
								Required: false,
								Value:    "npm_,KUBERNETES_,API_PORT",
							},
							&cli.StringFlag{
								Name:     "exclude",
								Usage:    "List (csv) of specific env vars to exclude values from display. Set to \"\" to remove any exclusions.",
								Required: false,
								Value:    "PATH,SHLVL,HOSTNAME",
							},
						},
						Action: cmd.ViewEnvDiff,
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
