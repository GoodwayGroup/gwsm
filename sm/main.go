package sm

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gwsm/lib"
	"os"
	"strings"
)

/*
   NOTE: Refactor of https://github.com/cyberark/summon-aws-secrets/blob/master/main.go
   This was needed to have the command return the byte stream rather than have it write
   to STDOUT
*/
func RetrieveSecret(variableName string) (secretBytes []byte) {
	sess, _ := lib.GetAWSSession()

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.

	svc := secretsmanager.New(sess)

	// Check if key has been specified
	arguments := strings.SplitN(variableName, "#", 2)

	secretName := arguments[0]
	var keyName string

	if len(arguments) > 1 {
		keyName = arguments[1]
	}

	// Get secret value
	req, resp := svc.GetSecretValueRequest(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})

	err := req.Send()
	if err != nil { // resp is now filled
		printAndExit(err)
	}

	if resp.SecretString != nil {
		secretBytes = []byte(*resp.SecretString)
	} else {
		secretBytes = resp.SecretBinary
	}

	if keyName != "" {
		secretBytes, err = getValueByKey(keyName, secretBytes)
		if err != nil {
			printAndExit(err)
		}
	}

	return
}

// TODO: Refactor to use common error return pattern for cli.Error
func printAndExit(err error) {
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}

func getValueByKey(keyName string, secretBytes []byte) (secret []byte, err error) {
	var secrets map[string]interface{}
	var secretValue string

	if err := json.Unmarshal(secretBytes, &secrets); err != nil {
		return nil, err
	}

	secretValue = fmt.Sprint(secrets[keyName])

	return []byte(secretValue), nil
}
