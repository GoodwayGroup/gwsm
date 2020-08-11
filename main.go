package main

import (
	"fmt"
	"github.com/GoodwayGroup/gwsm/cmd"
	"github.com/GoodwayGroup/gwsm/info"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

var version string

func main() {
	app := &cli.App{
		Name:     info.AppName,
		Version:  version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Derek Smith",
				Email: "dsmith@goodwaygroup.com",
			},
			{
				Name: info.AppRepoOwner,
			},
		},
		Copyright:            "(c) 2020 Goodway Group",
		HelpName:             info.AppName,
		Usage:                "interact with config map and secret manager variables",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "sm",
				Aliases: []string{"secretsmanager"},
				Usage:   "Secrets Manager commands w/ interactive interface",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "binary",
						Aliases:     []string{"b"},
						Usage:       "get the SecretBinary value",
						DefaultText: info.SecretBinaryHelp,
					},
				},
				Subcommands: []*cli.Command{
					{
						// list-secrets
						Name:   "list",
						Usage:  "display table of all secrets with meta data",
						Action: cmd.SMListSecrets,
					},
					{
						// describe-secret
						Name:  "describe",
						Usage: "print description of secret to `STDOUT`",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "secret-id",
								Aliases: []string{"s"},
								Usage:   "Specific Secret to describe, will bypass select/search",
							},
						},
						Action: cmd.SMDescribeSecret,
					},
					{
						// get-secret-value
						Name:    "get",
						Aliases: []string{"view"},
						Usage:   "select from list or pass in specific secret",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "secret-id",
								Aliases: []string{"s"},
								Usage:   "Specific Secret to view, will bypass select/search",
							},
						},
						Action: cmd.SMViewSecret,
					},
					{
						Name:    "edit",
						Aliases: []string{"e"},
						Usage:   "interactive edit of a secret String Value",
						// TODO: add UsageText
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "secret-id",
								Aliases: []string{"s"},
								Usage:   "Specific Secret to edit, will bypass select/search",
							},
						},
						Action: cmd.SMEditSecret,
					},
					{
						// create-secret
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "create new secret in Secrets Manager",
						// TODO: add UsageText
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secret-id",
								Aliases:  []string{"s"},
								Usage:    "Secret name",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "value",
								Aliases: []string{"v"},
								Usage:   "Secret Value. Will store as a string, unless binary flag is set.",
							},
							&cli.BoolFlag{
								Name:    "interactive",
								Aliases: []string{"i"},
								Usage:   "Open interactive editor to create secret value.",
							},
							&cli.StringFlag{
								// TODO: add description feature
								Name:    "description",
								Aliases: []string{"desc"},
								Usage:   "Additional description text.",
							},
							&cli.StringFlag{
								// TODO: add tags feature
								Name:  "tags",
								Usage: "key=value tags (CSV list)",
							},
						},
						Action: cmd.SMCreateSecret,
					},
					{
						// put-secret-value
						Name:  "put",
						Usage: "non-interactive update to a specific secret",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secret-id",
								Aliases:  []string{"s"},
								Usage:    "Secret name",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "value",
								Aliases: []string{"v"},
								Usage:   "Secret Value. Will store as a string, unless binary flag is set.",
							},
							&cli.BoolFlag{
								Name:    "interactive",
								Aliases: []string{"i"},
								Usage:   "Open interactive editor to create secret value.",
							},
							&cli.StringFlag{
								// TODO: add description feature
								Name:    "description",
								Aliases: []string{"desc"},
								Usage:   "Additional description text.",
							},
							&cli.StringFlag{
								// TODO: add tags feature
								Name:  "tags",
								Usage: "key=value tags (CSV list)",
							},
						},
						// TODO: Flag for use of binary
						Action: cmd.SMPutSecret,
					},
					{
						Name:    "delete",
						Aliases: []string{"del"},
						Usage:   "delete a specific secret",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secret-id",
								Aliases:  []string{"s"},
								Usage:    "Specific Secret to delete",
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Usage:   "Bypass recovery window (30 days) and immediately delete Secret.",
							},
						},
						Action: cmd.SMDeleteSecret,
					},
				},
			},
			{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "Commands to interact with environment variables, both local and on cluster.",
				Subcommands: []*cli.Command{
					{
						Name:    "diff",
						Aliases: []string{"d"},
						Usage:   "Print out detailed diff reports comparing local and running Pod",
						Subcommands: []*cli.Command{
							{
								Name:      "namespace",
								Aliases:   []string{"ns"},
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
										Usage:    "Kube Namespace to list Pods from for inspection",
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
										Usage:    "List of prefixes (csv) used to filter values from display. Set to `\"\"` to remove any filters.",
										Required: false,
										Value:    "npm_,KUBERNETES_,API_PORT",
									},
									&cli.StringFlag{
										Name:     "exclude",
										Usage:    "List (csv) of specific env vars to exclude values from display. Set to `\"\"` to remove any exclusions.",
										Required: false,
										Value:    "PATH,SHLVL,HOSTNAME",
									},
									&cli.StringFlag{
										Name:  "secret-suffix",
										Usage: "Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup",
										Value: "_NAME",
									},
								},
								Action: cmd.ViewEnvDiff,
							},
							{
								Name:      "ansible",
								Aliases:   []string{"legacy"},
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
										Usage:    "Kube Namespace list Pods from for inspection",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "dotenv",
										Usage:    "Path to `.env` file on Pod",
										Required: false,
										Value:    "$PWD/.env",
									},
								},
								Action: cmd.ViewAnsibleEnvDiff,
							},
						},
					},
					{
						Name:    "view",
						Aliases: []string{"v"},
						Usage:   "View configured environment for either local or running on a Pod",
						Subcommands: []*cli.Command{
							{
								Name:      "configmap",
								Aliases:   []string{"c"},
								Usage:     "View env values based on local settings in a ConfigMap and secrets.yml",
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
									&cli.StringFlag{
										Name:  "secret-suffix",
										Usage: "Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup",
										Value: "_NAME",
									},
								},
								Action: cmd.ViewLocalEnv,
							},
							{
								Name:      "ansible",
								Aliases:   []string{"legacy"},
								Usage:     "View env values from ansible-vault encrypted Secret file.",
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
							{
								Name:      "namespace",
								Aliases:   []string{"ns"},
								Usage:     "Interact with env on a running Pod within a Namespace",
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
				Name:  "install-manpage",
				Usage: "Generate and install man page",
				Action: func(c *cli.Context) error {
					mp, _ := info.ToMan(c.App)
					err := ioutil.WriteFile("/usr/local/share/man/man8/gwsm.8", []byte(mp), 0644)
					if err != nil {
						return cli.NewExitError(fmt.Sprintf("Unable to install man page: %e", err), 2)
					}
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version info",
				Action: func(c *cli.Context) error {
					fmt.Printf("%s %s (%s/%s)\n", info.AppName, version, runtime.GOOS, runtime.GOARCH)
					return nil
				},
			},
		},
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := info.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := info.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
