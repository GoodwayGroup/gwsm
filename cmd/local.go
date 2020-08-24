package cmd

import (
	"fmt"
	"github.com/GoodwayGroup/gwsm/env"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strings"
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
			l := kemba.PickColor(group)
			t.SetTitle(fmt.Sprintf("From secret: %s", l.Sprintf(group)))
		}
		t.AppendHeader(table.Row{"Key", "Value"})

		keys := make([]string, 0, len(values))
		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			var value string
			if strings.HasSuffix(key, c.String("secret-suffix")) {
				l := kemba.PickColor(values[key])
				value = l.Sprint(values[key])
				key = l.Sprint(key)
			} else {
				value = values[key]
			}

			t.AppendRow([]interface{}{
				key,
				value,
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
