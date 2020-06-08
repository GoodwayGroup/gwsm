package env

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"gwsm/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func shouldProcessLine(c *cli.Context, ln string) bool {
	prefixes := strings.Split(c.String("filter-prefix"), ",")
	if len(prefixes) > 0 {
		for _, prefix := range prefixes {
			if strings.HasPrefix(ln, prefix) {
				return false
			}
		}
	}

	exclusions := strings.Split(c.String("exclude"), ",")
	if len(exclusions) > 0 {
		result := strings.SplitN(ln, "=", 2)
		if len(result) > 1 {
			return !containsString(exclusions, result[0])
		}
		return false
	}

	return true
}

func promptForPod(err error, clientset *kubernetes.Clientset, namespace string) (string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		// TODO: Consolidate logger
		fmt.Println(aurora.Red("✖ Failed to get pods"))
		return "", err
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
		// TODO: Consolidate logger
		fmt.Println(aurora.Sprintf(aurora.Red("✖ Prompt failed %e"), err))
		return "", err
	}

	return result, nil
}

func GetEnvFromPodProcess(c *cli.Context) (envMap map[string]string, err error) {
	err, clientset := kube.GetClient()
	if err != nil {
		// TODO: Consolidate logger
		fmt.Println(aurora.Red("✖ Failed to get kube client"))
		return nil, err
	}

	namespace := c.String("namespace")
	result, err := promptForPod(err, clientset, namespace)
	if err != nil {
		return nil, err
	}

	fmt.Println("")
	cmd := []string{"/bin/sh", "-c", fmt.Sprintf("strings /proc/$(ps faux | grep %s | tail -1 | awk '{print $2}')/environ", c.String("cmd"))}
	stdOut, _, err := kube.ExecCommandInContainerWithFullOutput(clientset, namespace, result, cmd)
	if err != nil {
		// TODO: Consolidate logger
		fmt.Println(aurora.Sprintf(aurora.Red("✖ Failed to execute command `%s` on pod %s with error: %e"), cmd, result, err))

		return nil, err
	}

	envMap = make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(stdOut))
	for scanner.Scan() {
		ln := scanner.Text()

		if shouldProcessLine(c, ln) {
			result := strings.SplitN(ln, "=", 2)
			if len(result) > 1 {
				envMap[result[0]] = fmt.Sprint(result[1])
			}
		}
	}

	return
}

func GetLegacyEnvFromPodProcess(c *cli.Context) (envMap map[string]string, err error) {
	err, clientset := kube.GetClient()
	if err != nil {
		// TODO: Consolidate logger
		fmt.Println(aurora.Red("✖ Failed to get kube client"))
		return nil, err
	}

	namespace := c.String("namespace")
	result, err := promptForPod(err, clientset, namespace)
	if err != nil {
		return nil, err
	}

	fmt.Println("")
	cmd := []string{"/bin/sh", "-c", fmt.Sprintf("cat %s", c.String("dotenv"))}
	stdOut, _, err := kube.ExecCommandInContainerWithFullOutput(clientset, namespace, result, cmd)
	if err != nil {
		// TODO: Consolidate logger
		fmt.Println(aurora.Sprintf(aurora.Red("✖ \"Failed to execute command `%s` on pod %s with error: %e\""), cmd, result, err))
		return nil, err
	}

	envMap = make(map[string]string)
	for _, envStr := range strings.Split(stdOut, "\n") {
		parts := strings.SplitN(envStr, "=", 2)
		if len(parts) > 1 {
			envMap[parts[0]] = parts[1]
		}
	}

	return
}
