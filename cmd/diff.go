package cmd

import (
	"bytes"
	"fmt"
	"github.com/GoodwayGroup/gwsm/env"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora/v3"
	"github.com/r3labs/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strings"
)

var (
	kl = kemba.New("gwsm:diff")
)

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getPaintedValue(group string, v int) string {
	if v == 0 {
		return fmt.Sprint(aurora.Gray(15, v))
	}
	switch group {
	case "no_change":
		return fmt.Sprint(aurora.Blue(v))
	case "create":
		return fmt.Sprint(aurora.Green(v))
	case "update":
		return fmt.Sprint(aurora.Yellow(v))
	case "delete":
		return fmt.Sprint(aurora.Red(v))
	default:
		return fmt.Sprintf("%d", v)
	}
}

func addRowToSummary(group string, found int, t table.Writer) {
	val := getPaintedValue(group, found)
	switch group {
	case "no_change":
		t.AppendRow([]interface{}{
			aurora.Blue("NC"),
			val,
			"No Change",
			"Local ENV values that are the same on the Pod",
		})
	case "create":
		t.AppendRow([]interface{}{
			aurora.Green("N"),
			val,
			"New",
			"Local ENV values that are NOT found on the selected Pod",
		})
	case "update":
		t.AppendRow([]interface{}{
			aurora.Yellow("U"),
			val,
			"Update",
			"Local ENV values that are DIFFERENT than the values found on the selected Pod",
		})
	case "delete":
		b := `ENV values that were found on the selected Pod, but are not found locally. 
This could be caused by system level ENV values (such as CWD) or it could 
indicate that a value is MISSING from the local ENV.`
		t.AppendRow([]interface{}{
			aurora.Red("PD"),
			val,
			"Possible Delete",
			b,
		})
	default:
		// This should not be reached.
		panic(fmt.Sprintf("Uknown group type: %s", group))
	}
}

func addRowToTable(change diff.Change, t table.Writer) {
	switch change.Type {
	case "no_change":
		// This means that the value is not contained in the Pod environment and will be added.
		t.AppendRow([]interface{}{
			aurora.Blue("NC"),
			change.Path[0],
			change.From,
		})
	case "create":
		// This means that the value is not contained in the Pod environment and will be added.
		t.AppendRow([]interface{}{
			aurora.Green("N"),
			change.Path[0],
			change.To,
		})
	case "update":
		// This denotes that there is a change in the local value compared to that on the Pod.

		dmp := diffmatchpatch.New()
		//
		diffs := dmp.DiffMain(change.From.(string), change.To.(string), false)

		t.AppendRow([]interface{}{
			aurora.Yellow("U"),
			change.Path[0],
			fmt.Sprint(aurora.Gray(15, change.From)),
		})
		t.AppendRow([]interface{}{
			"",
			"",
			printUpdateDiff(diffs),
		})
		t.AppendSeparator()
	case "delete":
		// This indicates that the value is present on the Pod, but not in the local env.
		// TODO: Format the row with color based on whether it is a system variable or not
		t.AppendRow([]interface{}{
			aurora.Red("PD"),
			change.Path[0],
			change.From,
		})
	default:
		// This should not be reached.
		panic(fmt.Sprintf("Uknown change type: %s", change.Type))
	}
}

func printUpdateDiff(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, d := range diffs {
		text := d.Text

		switch d.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString(aurora.Sprintf(aurora.Yellow(text)))
		case diffmatchpatch.DiffEqual:
			_, _ = buff.WriteString(text)
		}
	}

	return buff.String()
}

func printOutDiff(changelog diff.Changelog, envMapLocal map[string]string) {
	diffGroups := groupDiffWithLocal(changelog, envMapLocal)

	groups := []string{"no_change", "create", "update", "delete"}

	tSum := table.NewWriter()
	tSum.SetOutputMirror(os.Stdout)
	tSum.SetStyle(table.StyleLight)
	tSum.AppendHeader(table.Row{"Status", "Count", "Name", "Description"})

	for _, group := range groups {
		addRowToSummary(group, len(diffGroups[group]), tSum)
	}

	block := strings.Repeat("-", 79)
	fmt.Printf("%s\n%s\n%s\n", block, "> Summary Overview", block)
	tSum.Render()
	fmt.Println("")

	for _, group := range groups {
		changes := len(diffGroups[group])

		if changes > 0 {
			fmt.Printf("%s\n> %s details\n%s\n", block, group, block)
			sort.Slice(diffGroups[group], func(i int, j int) bool {
				return diffGroups[group][i].Path[0] < diffGroups[group][j].Path[0]
			})

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleLight)
			t.AppendHeader(table.Row{"Status", "Key", "Value"})

			for _, change := range diffGroups[group] {
				addRowToTable(change, t)
			}

			t.Render()
			fmt.Println("")
		}
	}
}

func groupDiffWithLocal(changelog diff.Changelog, envMapLocal map[string]string) map[string][]diff.Change {
	diffGroups := map[string][]diff.Change{
		"create":    {},
		"delete":    {},
		"update":    {},
		"no_change": {},
	}

	var diffKeys []string
	for _, change := range changelog {
		diffGroups[change.Type] = append(diffGroups[change.Type], change)
		diffKeys = append(diffKeys, change.Path[0])
	}

	for k, v := range envMapLocal {
		if !containsString(diffKeys, k) {
			diffGroups["no_change"] = append(diffGroups["no_change"], diff.Change{
				Type: "no_change",
				Path: []string{k},
				From: v,
				To:   "",
			})
		}
	}
	return diffGroups
}

// ViewEnvDiff will print the diff for a parsed local ConfigMap file and
// retrieved JSON blobs from AWS Secrets Manager with the environment
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
	var envMapPod map[string]string
	envMapPod, err = env.GetEnvFromPodProcess(c)
	if err != nil {
		return err
	}

	// Compares as if the Local env is being applied to the Pod env.
	var changelog diff.Changelog
	changelog, err = diff.Diff(envMapPod, envMapLocal)
	if err != nil {
		return err
	}
	kl.Log(changelog)

	printOutDiff(changelog, envMapLocal)

	return nil
}

// ViewAnsibleEnvDiff will print the diff for a decrypted local Kube Secrets file
// with the environment pulled from the dotenv file for a Pod within a supplied
// NameSpace.
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
	var envMapPod map[string]string
	envMapPod, err = env.GetLegacyEnvFromPodProcess(c)
	if err != nil {
		return err
	}

	// Compares as if the Local env is being applied to the Pod env.
	var changelog diff.Changelog
	changelog, err = diff.Diff(envMapPod, envMapLocal)
	if err != nil {
		return err
	}

	printOutDiff(changelog, envMapLocal)

	return nil
}
