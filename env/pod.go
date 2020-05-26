package env

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"gwsm/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func GetEnvFromPodProcess(c *cli.Context) (envMap map[string]string, err error) {
	// TODO: Handle error
	err, clientset := kube.GetClient()
	if err != nil {
		fmt.Println("Failed to get kube client:", err)
		return nil, err
	}

	namespace := c.String("namespace")
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("Failed to get pods:", err)
		return nil, err
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
		return nil, err
	}

	fmt.Println("")
	cmd := []string{"/bin/sh", "-c", fmt.Sprintf("strings /proc/$(ps faux | grep %s | tail -1 | awk '{print $2}')/environ", c.String("cmd"))}
	stdOut, _, err := kube.ExecCommandInContainerWithFullOutput(clientset, namespace, result, cmd)
	if err != nil {
		fmt.Printf("Failed to execute command `%s` on pod %s with error: %e", cmd, result, err)
		return nil, err
	}

	envMap = make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(stdOut))
	for scanner.Scan() {
		ln := scanner.Text()
		if !strings.HasPrefix(scanner.Text(), c.String("filter")) {
			result := strings.SplitN(ln, "=", 2)
			if len(result) > 1 {
				envMap[result[0]] = fmt.Sprint(result[1])
			}
		}
	}

	return
}
