package env

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/pbthorste/avtool"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v3"
	"gwsm/lib"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

func GetEnvFromAnsibleVault(c *cli.Context) (string, error) {
	pw, err := retrieveVaultPassword(c.String("vault-password-file"))
	if err != nil {
		return "", err
	}

	vf := c.String("encrypted-env-file")
	result, err := avtool.DecryptFile(vf, pw)
	if err != nil {
		if strings.Compare(err.Error(), "ERROR: runtime error: index out of range") == 0 {
			return "", cli.NewExitError("input is not a vault encrypted "+vf+" is not a vault encrypted file for "+vf, 2)
		}
		return "", cli.NewExitError(err, 1)
	}

	var kubeSecret lib.KubeSecret
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
			return "", errors.New("ERROR: vault-password-file, could not find: " + vaultPasswordFile)
		}
		pw, err := ioutil.ReadFile(vaultPasswordFile)
		if err != nil {
			return "", errors.New("ERROR: vault-password-file, " + err.Error())
		}
		return strings.TrimSpace(string(pw)), nil
	}

	return readVaultPassword()
}

func readVaultPassword() (password string, err error) {
	println("Vault password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		err = errors.New("ERROR: could not input password, " + err.Error())
		return
	}
	password = string(bytePassword)
	return
}
