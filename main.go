package main

import (
	"fmt"
	"github.com/GoodwayGroup/gwsm/cmd"
	"github.com/GoodwayGroup/gwsm/info"
	"github.com/clok/cdocs"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"time"
)

var version string

func main() {
	// Generate the install-manpage command
	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: info.AppName,
	})
	if err != nil {
		log.Fatal(err)
	}

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
						Action: cmd.ListSecrets,
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
						Action: cmd.DescribeSecret,
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
						Action: cmd.ViewSecret,
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
							// TODO: add flag for passing version stage
						},
						Action: cmd.EditSecret,
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
								Usage:   "Open interactive editor to create secret value. If no 'value' is provided, an editor will be opened by default.",
							},
							&cli.StringFlag{
								Name:    "description",
								Aliases: []string{"d"},
								Usage:   "Additional description text.",
							},
							&cli.StringFlag{
								Name:    "tags",
								Aliases: []string{"t"},
								Usage:   "key=value tags (CSV list)",
							},
						},
						Action: cmd.CreateSecret,
					},
					{
						// put-secret-value
						Name:  "put",
						Usage: "non-interactive update to a specific secret",
						UsageText: `
Stores a new encrypted secret value in the specified secret. To do this, the 
operation creates a new version and attaches it to the secret. The version 
can contain a new SecretString value or a new SecretBinary value.

This will put the value to AWSCURRENT and retain one previous version 
with AWSPREVIOUS.
`,
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
								Usage:   "Override and open interactive editor to verify and modify the new secret value.",
							},
							// TODO: add flag for passing version stage
						},
						Action: cmd.PutSecret,
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
						Action: cmd.DeleteSecret,
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
			im,
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
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
