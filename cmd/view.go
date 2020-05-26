package cmd

import (
	"fmt"
	"github.com/r3labs/diff"
	"github.com/urfave/cli/v2"
	"gwsm/env"
	"sort"
)

// Print the resulting environment for a set of local ConfigMap and Summon secrets.yml file.
func ViewLocalEnv(c *cli.Context) error {
	groupedValues, err := env.GetGroupedLocalEnv(c)
	if err != nil {
		return err
	}

	var envValues []string
	for group, values := range groupedValues {
		if group == "local" {
			envValues = append(envValues, fmt.Sprintf("\n# from ConfigMap"))
		} else {
			envValues = append(envValues, fmt.Sprintf("\n# from secret: %s", group))
		}

		for key, value := range values {
			envValues = append(envValues, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// dump full env to screen
	for _, v := range envValues {
		fmt.Println(v)
	}

	return nil
}

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

// Print the diff for a parsed local ConfigMap file and retrieved JSON blobs from AWS Secrets Manager with the environment
// for a given process on a Pod within a supplied NameSpace.
func ViewEnvDiff(c *cli.Context) error {
	// Get local envMap
	groupedValues, err := env.GetGroupedLocalEnv(c)
	if err != nil {
		return err
	}
	envMapLocal := make(map[string]string)
	for _, group := range groupedValues {
		for k, v := range group {
			envMapLocal[k] = v
		}
	}

	// Get envMap from Pod
	envMapPod, err := env.GetEnvFromPodProcess(c)
	if err != nil {
		return err
	}

	// Compares as if the Local env is being applied to the Pod env.
	changelog, err := diff.Diff(envMapPod, envMapLocal)
	for _, change := range changelog {
		switch change.Type {
		case "create":
			// This means that the value is not contained in the Pod environment and will be added.
			fmt.Printf("NEW KEY: %s\n\tVALUE: %s\n", change.Path[0], change.To)
		case "update":
			// This denotes that there is a change in the local value compared to that on the Pod.
			fmt.Printf("UPDATED KEY: %s\n\t%s -> %s\n", change.Path[0], change.From, change.To)
		case "delete":
			// This indicates that the value is present on the Pod, but not in the local env.
			fmt.Printf("DELETED KEY: %s\n", change.Path[0])
		default:
			// This should not be reached.
			fmt.Println(change)
		}
	}

	return nil
}
