package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/r3labs/diff"
	"github.com/urfave/cli/v2"
	"gwsm/env"
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
		t.SetStyle(table.StyleColoredDark)
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

	diffGroups := map[string][]diff.Change{
		"create": []diff.Change{},
		"delete": []diff.Change{},
		"update": []diff.Change{},
	}

	for _, change := range changelog {
		diffGroups[change.Type] = append(diffGroups[change.Type], change)
	}

	for _, group := range []string{"create", "update", "delete"} {
		changes := len(diffGroups[group])

		if changes > 0 {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleColoredDark)

			var title string
			switch group {
			case "create":
				title = "NEW: Env values that will be ADDED to the Pod"
			case "update":
				title = "UPDATE: Env values that will be UPDATED on the Pod"
			case "delete":
				title = "POSSIBLE DELETE: Env values that are present on the Pod, but not found in the local Env"
			default:
				// This should not be reached.
				panic(fmt.Sprintf("Uknown group type: %s", group))
			}

			t.SetTitle(title)
			t.AppendHeader(table.Row{"Key", "Value"})

			fmt.Printf("Found %d differences for group `%s`\n", changes, group)
			sort.Slice(diffGroups[group], func(i int, j int) bool {
				return diffGroups[group][i].Path[0] < diffGroups[group][j].Path[0]
			})
			for _, change := range diffGroups[group] {
				switch change.Type {
				case "create":
					// This means that the value is not contained in the Pod environment and will be added.
					// fmt.Printf("NEW KEY: %s=%s\n", change.Path[0], change.To)
					t.AppendRow([]interface{}{
						change.Path[0],
						change.To,
					})
				case "update":
					// This denotes that there is a change in the local value compared to that on the Pod.
					// fmt.Printf("UPDATED KEY: %s\n\t%s -> %s\n", change.Path[0], change.From, change.To)
					t.AppendRow([]interface{}{
						change.Path[0],
						fmt.Sprintf("%s -> %s", change.From, change.To),
					})
				case "delete":
					// This indicates that the value is present on the Pod, but not in the local env.
					// fmt.Printf("DELETED KEY: %s=%s\n", change.Path[0], change.From)
					t.AppendRow([]interface{}{
						change.Path[0],
						change.From,
					})
				default:
					// This should not be reached.
					panic(fmt.Sprintf("Uknown change type: %s", change.Type))
				}
			}
			t.Render()
		} else {
			fmt.Printf("No diff found for group `%s`\n", group)
		}
		fmt.Println("")
	}

	return nil
}
