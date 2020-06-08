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
	if err != nil {
		// TODO: Refactor to use common error return pattern for cli.Error
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
			// TODO: Refactor to use common error return pattern for cli.Error
			printAndExit(err)
		}
	}

	return
}

// ListSecrets will retrieval ALL secrets via pagination of 100 per page. It will
// return once all pages have been processed.
func ListSecrets() (secrets []secretsmanager.SecretListEntry, err error) {
	sess, _ := lib.GetAWSSession()
	svc := secretsmanager.New(sess)

	// Get all secret names
	err = svc.ListSecretsPages(&secretsmanager.ListSecretsInput{
		MaxResults: aws.Int64(100),
	},
		func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
			for _, v := range page.SecretList {
				secrets = append(secrets, *v)
			}
			return !lastPage
		})
	if err != nil {
		return nil, err
	}

	return
}

// GetSecret will retrieve a specific secret by Name (id)
func GetSecret(id string) (secret *secretsmanager.GetSecretValueOutput, err error) {
	sess, _ := lib.GetAWSSession()
	svc := secretsmanager.New(sess)

	secret, err = svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return
}

// PutSecretString will put an updated SecretString value to a specific secret by Name (id)
func PutSecretString(id string, data string) (secret *secretsmanager.PutSecretValueOutput, err error) {
	sess, _ := lib.GetAWSSession()
	svc := secretsmanager.New(sess)

	secret, err = svc.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretString: aws.String(data),
		SecretId:     aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return
}

// PutSecretBinary will put an updated SecretBinary value to a specific secret by Name (id)
func PutSecretBinary(id string, data []byte) (secret *secretsmanager.PutSecretValueOutput, err error) {
	sess, _ := lib.GetAWSSession()
	svc := secretsmanager.New(sess)

	secret, err = svc.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretBinary: data,
		SecretId:     aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return
}

// CreateSecretString will create a new SecretString value to a specific secret by Name (id)
func CreateSecretString(id string, data string) (secret *secretsmanager.CreateSecretOutput, err error) {
	sess, _ := lib.GetAWSSession()
	svc := secretsmanager.New(sess)

	secret, err = svc.CreateSecret(&secretsmanager.CreateSecretInput{
		SecretString: aws.String(data),
		Name:         aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return
}

// CreateSecretBinary will create a new SecretBinary value to a specific secret by Name (id)
func CreateSecretBinary(id string, data []byte) (secret *secretsmanager.CreateSecretOutput, err error) {
	sess, _ := lib.GetAWSSession()
	svc := secretsmanager.New(sess)

	secret, err = svc.CreateSecret(&secretsmanager.CreateSecretInput{
		SecretBinary: data,
		Name:         aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return
}

// DescribeSecret retrieves the describe data for a specific secret by Name (id)
func DescribeSecret(id string) (secret *secretsmanager.DescribeSecretOutput, err error) {
	sess, _ := lib.GetAWSSession()
	svc := secretsmanager.New(sess)

	// Get all secret names
	secret, err = svc.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: aws.String(id),
	})

	if err != nil {
		return nil, err
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
