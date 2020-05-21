package sm

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"os"
	"strings"
)

/*
    NOTE: Refactor of https://github.com/cyberark/summon-aws-secrets/blob/master/main.go
    This was needed to have the command return the byte stream rather than have it write
    to STDOUT
 */
func RetrieveSecret(variableName string) (secretBytes []byte) {
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		printAndExit(err)
	}

	// AWS Go SDK does not currently support automatic fetching of region from ec2metadata.
	// If the region could not be found in an environment variable or a shared config file,
	// create metaSession to fetch the ec2 instance region and pass to the regular Session.
	if *sess.Config.Region == "" {
		metaSession, err := session.NewSession()
		if err != nil {
			printAndExit(err)
		}

		metaClient := ec2metadata.New(metaSession)
		// If running on an EC2 instance, the metaClient will be available and we can set the region to match the instance
		// If not on an EC2 instance, the region will remain blank and AWS returns a "MissingRegion: ..." error
		if metaClient.Available() {
			if region, err := metaClient.Region(); err == nil {
				sess.Config.Region = aws.String(region)
			} else {
				printAndExit(err)
			}
		}
	}

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

	err = req.Send()
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
