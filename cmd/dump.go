package cmd

import (
	"fmt"
	"github.com/GoodwayGroup/gwsm/env"
	"github.com/clok/kemba"
	"github.com/urfave/cli/v2"
	"sort"
	"strings"
)

// Print the resulting environment for a set of local ConfigMap and Summon secrets.yml file.
func DumpLocalEnv(c *cli.Context) error {
	groupedValues, err := env.GetGroupedLocalEnv(c)
	if err != nil {
		return err
	}

	k := kemba.New("DumpLocalEnv")
	k.Log(groupedValues)

	for scope, values := range groupedValues {
		fmt.Printf("# From: %s\n", scope)
		keys := make([]string, 0, len(values))
		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			fmt.Printf("%s=%s\n", key, values[key])
		}
		fmt.Println("")
	}

	return nil
}

func DumpAnsibleEncryptedEnv(c *cli.Context) error {
	dataStr, err := env.GetEnvFromAnsibleVault(c)
	if err != nil {
		return err
	}

	fmt.Println(strings.TrimSuffix(dataStr, "\n"))
	return nil
}
