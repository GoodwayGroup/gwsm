package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cyberark/summon/secretsyml"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"gwsm/kube"
	"gwsm/sm"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"sort"
	"strings"
	"sync"
)

func ViewLocalEnv(c *cli.Context) error {
	yamlFile, err := ioutil.ReadFile(c.String("configmap-path"))
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return err
	}

	var yamlConfig ConfigMap
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	var env []string
	env = append(env, "# from ConfigMap")
	subs := make(map[string]string)
	for key, value := range yamlConfig.Data {
		if strings.HasSuffix(key, "_NAME") {
			subs[strings.ToLower(key)] = value
		} else {
			env = append(env, fmt.Sprintf("%s=%s", key, fmt.Sprint(value)))
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
		return nil
	}

	secrets, err := secretsyml.ParseFromFile(c.String("secrets-path"), "", subs)
	if err != nil {
		fmt.Printf("Error parsing Secrets file: %s\n", err)
	}

	results := make(chan Result, len(secrets))
	var wg sync.WaitGroup

	for _, secretGroup := range subs {
		wg.Add(1)
		go func(groupName string) {
			defer wg.Done()
			blob := sm.RetrieveSecret(groupName)
			var parsed map[string]interface{}
			// TODO: address err capture
			json.Unmarshal(blob, &parsed)
			results <- Result{groupName, parsed, nil}
		}(secretGroup)
	}
	wg.Wait()
	close(results)

	allSecrets := make(map[string]map[string]interface{})
	for secretJSON := range results {
		allSecrets[secretJSON.Name] = secretJSON.JSON
	}

	groupedValues := make(map[string]map[string]string)
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

	for group, values := range groupedValues {
		env = append(env, fmt.Sprintf("\n# from secret: %s", group))
		for key, value := range values {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// dump full env to screen
	for _, v := range env {
		fmt.Println(v)
	}

	return nil
}

func ViewNamespaceEnv(c *cli.Context) error {
	// TODO: Handle error
	_, clientset := kube.GetClient()

	namespace := c.String("namespace")
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.GetName())
	}

	prompt := promptui.Select{
		Label: "Select Pod",
		Items: podNames,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	fmt.Println("")
	cmd := []string{"/bin/sh", "-c", fmt.Sprintf("strings /proc/$(ps faux | grep %s | tail -1 | awk '{print $2}')/environ", c.String("cmd"))}
	stdOut, _, err := kube.ExecCommandInContainerWithFullOutput(clientset, namespace, result, cmd)
	if err != nil {
		panic(err)
	}

	var env []string
	scanner := bufio.NewScanner(strings.NewReader(stdOut))
	for scanner.Scan() {
		ln := scanner.Text()
		if !strings.HasPrefix(scanner.Text(), c.String("filter")) {
			env = append(env, ln)
		}
	}

	sort.Strings(env)
	for _, val := range env {
		fmt.Println(val)
	}

	return nil
}
