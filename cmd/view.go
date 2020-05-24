package cmd

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"gwsm/env"
	"gwsm/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"sort"
	"strings"
)

func ViewLocalEnv(c *cli.Context) error {
	groupedValues, err := env.GetGroupedLocalEnv(c)
	if err != nil {
		log.Fatalln("failed to get pods:", err)

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

func ViewEnvDiff(c *cli.Context) error {
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
