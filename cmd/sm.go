package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/TylerBrock/colorjson"
	"github.com/a8m/djson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"gwsm/lib"
	"gwsm/sm"
	"os"
	"sort"
)

// Limit the length of a string while also appending an ellipses.
func truncateString(str string, num int) string {
	short := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		short = str[0:num] + "..."
	}
	return short
}

// Helper method to either bypass and return the `secretName` passed in via CLI
// flag OR retrieve a list of all secrets to allow for a search select by the
// User.
func selectSecretNameFromList(c *cli.Context) (string, error) {
	secretName := c.String("secret-id")
	if secretName == "" {
		secrets, err := sm.ListSecrets()
		if err != nil {
			PrintWarn("Error retrieving list of secrets.")
			return "", err
		}

		secretNames := make([]string, 0, len(secrets))
		for _, secret := range secrets {
			secretNames = append(secretNames, aws.StringValue(secret.Name))
		}
		sort.Strings(secretNames)

		var qs = []*survey.Question{
			{
				Name: "secret",
				Prompt: &survey.Select{
					Message: "Choose a Secret to view:",
					Options: secretNames,
					Default: secretNames[0],
				},
			},
		}

		answers := struct {
			SecretName string `survey:"secret"`
		}{}

		// perform the questions
		err = survey.Ask(qs, &answers)
		if err != nil {
			return "", err
		}

		secretName = answers.SecretName
		fmt.Printf("Retrieving: %s\n", secretName)
	}
	return secretName, nil
}

func SMListSecrets(c *cli.Context) error {
	secrets, err := sm.ListSecrets()
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Name", "Updated", "Accessed", "Description"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Name", WidthMax: 120},
		{Name: "Updated", WidthMax: 10},
		{Name: "Accessed", WidthMax: 10},
		{
			Name:     "Description",
			WidthMax: 40,
		},
	})
	t.SortBy([]table.SortBy{
		{Name: "Name", Mode: table.Asc},
	})

	for _, secret := range secrets {
		lastdt := aws.TimeValue(secret.LastAccessedDate)
		updateddt := aws.TimeValue(secret.LastChangedDate)
		t.AppendRow([]interface{}{
			aws.StringValue(secret.Name),
			fmt.Sprintf("%d-%02d-%02d", updateddt.Year(), updateddt.Month(), updateddt.Day()),
			fmt.Sprintf("%d-%02d-%02d", lastdt.Year(), lastdt.Month(), lastdt.Day()),
			truncateString(aws.StringValue(secret.Description), 40),
		})
	}

	t.Render()

	return nil
}

func SMViewSecret(c *cli.Context) error {
	secretName, err := selectSecretNameFromList(c)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	secret, err := sm.GetSecret(secretName)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	if c.Bool("binary") {
		fmt.Println(string(secret.SecretBinary))
	} else {
		result, err := djson.Decode([]byte(aws.StringValue(secret.SecretString)))
		if err != nil {
			PrintWarn("stored string value is not valid JSON.")
			fmt.Println(secret.SecretString)
		} else {
			f := colorjson.NewFormatter()
			f.Indent = 4

			s, _ := f.Marshal(result)
			fmt.Println(string(s))
		}
	}

	return nil
}

func SMDescribeSecret(c *cli.Context) error {
	secretName, err := selectSecretNameFromList(c)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	secret, err := sm.DescribeSecret(secretName)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	fmt.Println(secret.String())

	return nil
}

func SMEditSecret(c *cli.Context) error {
	secretName, err := selectSecretNameFromList(c)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	secret, err := sm.GetSecret(secretName)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	var s []byte
	if c.Bool("binary") {
		s = secret.SecretBinary
	} else {
		result, err := djson.Decode([]byte(aws.StringValue(secret.SecretString)))
		if err != nil {
			PrintWarn("stored string value is not valid JSON.")
			s = []byte(aws.StringValue(secret.SecretString))
		} else {
			s, err = json.MarshalIndent(result, "", "    ")
			if err != nil {
				return cli.NewExitError(err, 2)
			}
		}
	}

	up, err := lib.GetInputFromEditor(s)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	// TODO: Check if new data is Valid JSON before save

	if c.Bool("binary") {
		_, err = sm.PutSecretBinary(secretName, up)
	} else {
		_, err = sm.PutSecretString(secretName, string(up))
	}

	if err != nil {
		return cli.NewExitError(err, 2)
	}

	PrintSuccess(fmt.Sprintf("%s successfully updated.", secretName))

	return nil
}

func SMCreateSecret(c *cli.Context) error {
	secretName := c.String("secret-id")
	exists, err := sm.CheckIfSecretExists(secretName)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	if exists {
		PrintWarn(fmt.Sprintf("'%s' already exists. Please use a different name.", secretName))
		return nil
	}

	interactive := c.Bool("interactive")
	var value []byte
	if c.String("value") == "" {
		// Assume interactive mode
		interactive = true
		value = []byte("{}")
	} else {
		value = []byte(c.String("value"))
	}

	var s []byte
	if interactive {
		result, err := djson.Decode(value)
		if err != nil {
			PrintWarn("value is not valid JSON.")
			s = value
		} else {
			s, err = json.MarshalIndent(result, "", "    ")
			if err != nil {
				return cli.NewExitError(err, 2)
			}
		}

		up, err := lib.GetInputFromEditor(s)
		if err != nil {
			return cli.NewExitError(err, 2)
		}
		s = up
	}

	var t string
	if c.Bool("binary") {
		t = "BinarySecret"
		_, err = sm.CreateSecretBinary(secretName, s)
	} else {
		t = "StringSecret"
		_, err = sm.CreateSecretString(secretName, string(s))
	}

	if err != nil {
		return cli.NewExitError(err, 2)
	}

	PrintSuccess(fmt.Sprintf("%s %s successfully created.", secretName, t))

	return nil
}

func SMPutSecret(c *cli.Context) error {
	secretName := c.String("secret-id")
	exists, err := sm.CheckIfSecretExists(secretName)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	if exists {
		PrintWarn(fmt.Sprintf("'%s' already exists. Please use a different name.", secretName))
		return nil
	}

	// TODO: Implement PutSecret
	return cli.NewExitError("Not yet implemented", 5)
}
