package env

import (
	"encoding/json"
	"fmt"
	"github.com/cyberark/summon/secretsyml"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"gwsm/lib"
	"gwsm/sm"
	"io/ioutil"
	"strings"
	"sync"
)

// Parse local ConfigMap file and retrieve JSON blobs from AWS Secrets Manager. Return a map of groups names to value blocks.
func GetGroupedLocalEnv(c *cli.Context) (groupedValues map[string]map[string]string, err error) {
	yamlFile, err := ioutil.ReadFile(c.String("configmap"))
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return nil, err
	}

	var yamlConfig lib.ConfigMap
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	groupedValues = make(map[string]map[string]string)
	subs := make(map[string]string)
	for key, value := range yamlConfig.Data {
		// TODO: Make Secrets suffix a CLI param w/ default
		if strings.HasSuffix(key, "_NAME") {
			subs[strings.ToLower(key)] = value
		} else {
			if len(groupedValues["local"]) > 0 {
				groupedValues["local"][key] = fmt.Sprint(value)
			} else {
				groupedValues["local"] = map[string]string{key: fmt.Sprint(value)}
			}
		}
	}

	cnt := len(subs)
	if cnt > 0 {
		fmt.Printf("Found %d AWS Secrets Groups:\n", cnt)
		for key, value := range subs {
			fmt.Printf("\t%s: %s\n", strings.ToUpper(key), value)
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

	results := make(chan lib.Result, len(secrets))
	var wg sync.WaitGroup

	for _, secretGroup := range subs {
		wg.Add(1)
		go func(groupName string) {
			defer wg.Done()
			blob := sm.RetrieveSecret(groupName)
			var parsed map[string]interface{}
			err := json.Unmarshal(blob, &parsed)
			results <- lib.Result{Name: groupName, JSON: parsed, Error: err}
		}(secretGroup)
	}
	wg.Wait()
	close(results)

	allSecrets := make(map[string]map[string]interface{})
	for secretJSON := range results {
		if secretJSON.Error != nil {
			fmt.Printf("Error with SM results JSON: %s\n", secretJSON.Error)
			return nil, secretJSON.Error
		}
		allSecrets[secretJSON.Name] = secretJSON.JSON
	}

	for envvar, mappings := range secrets {
		arguments := strings.SplitN(mappings.Path, "#", 2)

		secretName := arguments[0]
		var keyName string

		if len(arguments) > 1 {
			keyName = arguments[1]
		}
		if len(groupedValues[secretName]) > 0 {
			groupedValues[secretName][envvar] = fmt.Sprint(allSecrets[secretName][keyName])
		} else {
			groupedValues[secretName] = map[string]string{envvar: fmt.Sprint(allSecrets[secretName][keyName])}
		}
	}

	return
}
