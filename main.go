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
			&cli.Author{
				Name: info.AppRepoOwner,
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
				Name:  "s3",
				Usage: "simple S3 commands",
				Subcommands: []*cli.Command{
					{
						Name:      "get",
						Usage:     "[object path] [destination path]",
						UsageText: info.S3GetCommandHelp,
						Action:    cmd.S3Get,
					},
				},
			},
			{
				Name:    "local",
				Aliases: []string{"l"},
				Usage:   "Interact with local env files",
				Subcommands: []*cli.Command{
					{
						Name:      "view",
						Aliases:   []string{"v"},
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
						Name:      "ansible",
						Aliases:   []string{"legacy", "a"},
						Usage:     "View value from ansible-vault encrypted Kube Secret file.",
						UsageText: info.ViewAnsibleEncryptedEnvCommandHelpText,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "vault-password-file",
								Usage:    "vault password file `VAULT_PASSWORD_FILE`",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "encrypted-env-file",
								Aliases:  []string{"e"},
								Usage:    "Path to encrypted Kube Secret file",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "accessor",
								Aliases: []string{"a"},
								Usage:   "Accessor key to pull data out of Data block.",
								Value:   ".env",
							},
						},
						Action: cmd.ViewAnsibleEncryptedEnv,
					},
				},
			},
			{
				Name:    "namespace",
				Aliases: []string{"ns"},
				Usage:   "Interact with env on a running Pod within a Namespace",
				Subcommands: []*cli.Command{
					{
						Name:      "view",
						Aliases:   []string{"v"},
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
				},
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
			{
				Name:      "diff:legacy",
				Aliases:   []string{"diff:ansible"},
				Usage:     "View diff of local (ansible encrypted) vs. namespace",
				UsageText: info.ViewAnsibleEnvDiffCommandHelpText,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "vault-password-file",
						Usage:    "vault password file `VAULT_PASSWORD_FILE`",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "encrypted-env-file",
						Aliases:  []string{"e"},
						Usage:    "Path to encrypted Kube Secret file",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "accessor",
						Aliases: []string{"a"},
						Usage:   "Accessor key to pull data out of Data block.",
						Value:   ".env",
					},
					&cli.StringFlag{
						Name:     "namespace",
						Aliases:  []string{"n"},
						Usage:    "Kube Namespace list Pods from",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "dotenv",
						Usage:    "Path to .env file on Pod",
						Required: false,
						Value:    "$PWD/.env",
					},
				},
				Action: cmd.ViewAnsibleEnvDiff,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
