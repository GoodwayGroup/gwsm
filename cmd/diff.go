package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/logrusorgru/aurora"
	"github.com/r3labs/diff"
	"github.com/urfave/cli/v2"
	"gwsm/env"
	"os"
	"sort"
	"strings"
)

func getDescription(group string, found int) string {
	switch group {
	case "create":
		return fmt.Sprintf(
			"%s\n\n%s\n\nFound %s",
			aurora.Green("New"),
			"Local ENV values that are NOT found on the selected Pod",
			fmt.Sprint(aurora.Green(found)),
		)
	case "update":
		return fmt.Sprintf(
			"%s\n\n%s\n\nFound %s",
			aurora.Yellow("Updates"),
			"Local ENV values that are DIFFERENT than the values found on the selected Pod",
			fmt.Sprint(aurora.Yellow(found)),
		)
	case "delete":
		b := `ENV values that were found on the selected Pod, but are not found locally. 
This could be caused by system level ENV values (such as CWD) or it could 
indicate that a value is MISSING from the local ENV.`
		return fmt.Sprintf(
			"%s\n\n%s\n\nFound %s",
			aurora.Red("Possible Deletions"),
			b,
			fmt.Sprint(aurora.Red(found)),
		)
	default:
		// This should not be reached.
		panic(fmt.Sprintf("Uknown group type: %s", group))
	}
}

func addRowToTable(change diff.Change, t table.Writer) {
	switch change.Type {
	case "create":
		// This means that the value is not contained in the Pod environment and will be added.
		t.AppendRow([]interface{}{
			change.Path[0],
			change.To,
		})
	case "update":
		// This denotes that there is a change in the local value compared to that on the Pod.
		t.AppendRow([]interface{}{
			change.Path[0],
			fmt.Sprintf("%s -> %s", aurora.Yellow(change.From), aurora.Green(change.To)),
		})
	case "delete":
		// This indicates that the value is present on the Pod, but not in the local env.
		// TODO: Format the row with color based on whether it is a system variable or not
		t.AppendRow([]interface{}{
			change.Path[0],
			change.From,
		})
	default:
		// This should not be reached.
		panic(fmt.Sprintf("Uknown change type: %s", change.Type))
	}
}

func printOutDiff(changelog diff.Changelog) {
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
		block := strings.Repeat("-", 79)
		fmt.Printf("%s\n%s\n%s\n", block, getDescription(group, changes), block)

		if changes > 0 {
			sort.Slice(diffGroups[group], func(i int, j int) bool {
				return diffGroups[group][i].Path[0] < diffGroups[group][j].Path[0]
			})

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleLight)
			t.AppendHeader(table.Row{"Key", "Value"})

			for _, change := range diffGroups[group] {
				addRowToTable(change, t)
			}

			t.Render()
		}
		fmt.Println("")
	}
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

	printOutDiff(changelog)

	return nil
}

// Print the diff for a decrypted local Kube Secrets file with the environment pulled from the dotenv file for
// a Pod within a supplied NameSpace.
func ViewAnsibleEnvDiff(c *cli.Context) error {
	// Get local envMap
	dataStr, err := env.GetEnvFromAnsibleVault(c)
	if err != nil {
		return err
	}

	envMapLocal := make(map[string]string)
	for _, envStr := range strings.Split(dataStr, "\n") {
		parts := strings.SplitN(envStr, "=", 2)
		if len(parts) > 1 {
			envMapLocal[parts[0]] = parts[1]
		}
	}

	// Get envMap from Pod
	envMapPod, err := env.GetLegacyEnvFromPodProcess(c)
	if err != nil {
		return err
	}

	// Compares as if the Local env is being applied to the Pod env.
	changelog, err := diff.Diff(envMapPod, envMapLocal)

	printOutDiff(changelog)

	return nil
}