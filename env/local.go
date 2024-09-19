package env

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/GoodwayGroup/gwsm/sm"
	"github.com/clok/kemba"
	"github.com/cyberark/summon/pkg/secretsyml"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// ConfigMap is a simple version of the ConfigMap defined in kubernetes.
type ConfigMap struct {
	Data map[string]string
}

// Result is a common interface for returning metadata to a channel.
type Result struct {
	Name  string
	JSON  map[string]interface{}
	Error error
}

func addToGroupedValues(groupedValues map[string]map[string]string, group string, key string, value string) {
	if len(groupedValues[group]) > 0 {
		groupedValues[group][key] = fmt.Sprint(value)
	} else {
		groupedValues[group] = map[string]string{key: value}
	}
}

// GetGroupedLocalEnv will parse local ConfigMap file and retrieve JSON blobs
// from AWS Secrets Manager. Return a map of groups names to value blocks.
func GetGroupedLocalEnv(c *cli.Context) (groupedValues map[string]map[string]string, err error) {
	yamlFile, err := os.ReadFile(c.String("configmap"))
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return nil, err
	}

	var yamlConfig ConfigMap
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	groupedValues = make(map[string]map[string]string)
	subs := make(map[string]string)
	for key, value := range yamlConfig.Data {
		if strings.HasSuffix(key, c.String("secret-suffix")) {
			subs[strings.ToLower(key)] = value
		}
		addToGroupedValues(groupedValues, "local", key, fmt.Sprint(value))
	}

	cnt := len(subs)
	if cnt > 0 {
		fmt.Printf("Found %d AWS Secrets Groups:\n", cnt)
		for key, value := range subs {
			l := kemba.PickColor(value)
			fmt.Printf("\t%s: %s\n", l.Sprint(strings.ToUpper(key)), l.Sprint(value))
		}
		fmt.Println("")
	} else {
		fmt.Println("No AWS Secrets Groups found.")
		return
	}

	secrets, err := secretsyml.ParseFromFile(c.String("secrets"), "", subs)
	if err != nil {
		fmt.Printf("Error parsing Secrets file: %s\n", err)
	}

	results := make(chan Result, len(secrets))
	var wg sync.WaitGroup

	for _, secretGroup := range subs {
		wg.Add(1)
		go func(groupName string) {
			defer wg.Done()
			blob, err := sm.RetrieveSecret(groupName)
			if err != nil {
				results <- Result{Name: groupName, JSON: nil, Error: err}
				return
			}
			var parsed map[string]interface{}
			err = json.Unmarshal(blob, &parsed)
			results <- Result{Name: groupName, JSON: parsed, Error: err}
		}(secretGroup)
	}
	wg.Wait()
	close(results)

	allSecrets := make(map[string]map[string]interface{})
	for secretJSON := range results {
		if secretJSON.Error != nil {
			fmt.Printf("Error with SM results JSON: %s\n", secretJSON.Error)
		} else {
			allSecrets[secretJSON.Name] = secretJSON.JSON
		}
	}

	for envvar, mappings := range secrets {
		arguments := strings.SplitN(mappings.Path, "#", 2)

		secretName := arguments[0]
		var keyName string

		if len(arguments) > 1 {
			keyName = arguments[1]
		}
		addToGroupedValues(groupedValues, secretName, envvar, fmt.Sprint(allSecrets[secretName][keyName]))
	}

	return
}
