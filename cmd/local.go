package cmd

import (
	"fmt"
	"github.com/GoodwayGroup/gwsm/env"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
)

// Print the resulting environment for a set of local ConfigMap and Summon secrets.yml file.
func ViewLocalEnv(c *cli.Context) error {
	groupedValues, err := env.GetGroupedLocalEnv(c)
	if err != nil {
		return err
	}

	for group, values := range groupedValues {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		if group == "local" {
			t.SetTitle("From ConfigMap")
		} else {
			t.SetTitle(fmt.Sprintf("From secret: %s", group))
		}
		t.AppendHeader(table.Row{"Key", "Value"})

		keys := make([]string, 0, len(values))
		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			t.AppendRow([]interface{}{
				key,
				values[key],
			})
		}
		t.Render()
		fmt.Println("")
	}

	return nil
}

func ViewAnsibleEncryptedEnv(c *cli.Context) error {
	dataStr, err := env.GetEnvFromAnsibleVault(c)
	if err != nil {
		return err
	}

	fmt.Println(dataStr)
	return nil
}
