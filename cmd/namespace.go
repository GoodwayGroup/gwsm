package cmd

import (
	"fmt"
	"github.com/GoodwayGroup/gwsm/env"
	"github.com/urfave/cli/v2"
	"sort"
)

// Print the environment for a given process on a Pod within a supplied NameSpace.
func ViewNamespaceEnv(c *cli.Context) error {
	envMap, err := env.GetEnvFromPodProcess(c)
	if err != nil {
		return err
	}

	sortedEnv := make([]string, 0, len(envMap))
	for k, v := range envMap {
		sortedEnv = append(sortedEnv, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(sortedEnv)

	for _, ln := range sortedEnv {
		fmt.Println(ln)
	}

	return nil
}
