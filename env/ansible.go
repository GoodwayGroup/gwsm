package env

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/clok/avtool/v3"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// KubeSecret is a simple version of the KubeSecret defined in kubernetes.
type KubeSecret struct {
	Kind string
	Type string
	Data map[string]interface{}
}

func GetEnvFromAnsibleVault(c *cli.Context) (string, error) {
	println("Retrieving vault password")
	pw, err := retrieveVaultPassword(c.String("vault-password-file"))
	if err != nil {
		return "", err
	}

	pwBytes := []byte(pw)
	vf := c.String("encrypted-env-file")
	result, err := avtool.DecryptFile(&avtool.DecryptFileOptions{
		Filename: vf,
		Password: &pwBytes,
	})
	if err != nil {
		if strings.Compare(err.Error(), "ERROR: runtime error: index out of range") == 0 {
			return "", cli.Exit("input is not a vault encrypted "+vf+" is not a vault encrypted file for "+vf, 2)
		}
		return "", cli.Exit(err, 1)
	}

	var kubeSecret KubeSecret
	err = yaml.Unmarshal([]byte(result), &kubeSecret)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	data, err := base64.StdEncoding.DecodeString(kubeSecret.Data[c.String("accessor")].(string))
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	dataStr := string(data)
	dataStr = strings.TrimSuffix(dataStr, "\"")
	dataStr = strings.TrimPrefix(dataStr, "\"")

	return dataStr, nil
}

func retrieveVaultPassword(vaultPasswordFile string) (string, error) {
	if vaultPasswordFile != "" {
		if _, err := os.Stat(vaultPasswordFile); os.IsNotExist(err) {
			return "", fmt.Errorf("ERROR: vault-password-file, could not find: %s", vaultPasswordFile)
		}

		pw, err := os.ReadFile(vaultPasswordFile)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(pw)), nil
	}

	return readVaultPassword()
}

func readVaultPassword() (password string, err error) {
	println("Vault password: ")
	var bytePassword []byte
	// skipped for windows compatibility
	// nolint:unconvert
	bytePassword, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		err = errors.New("ERROR: could not input password, " + err.Error())
		return
	}
	password = string(bytePassword)
	return
}
