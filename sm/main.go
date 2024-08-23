package sm

import (
	"encoding/json"
	"fmt"
	"context"
	"log"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"strings"
)

func getValueByKey(keyName string, secretBytes []byte) (secret []byte, err error) {
	var secrets map[string]interface{}
	var secretValue string

	if err := json.Unmarshal(secretBytes, &secrets); err != nil {
		return nil, err
	}

	secretValue = fmt.Sprint(secrets[keyName])

	return []byte(secretValue), nil
}

// RetrieveSecret will pull the AWS Secrets Manager value and parse out the specific value needed.
//
// NOTE: Refactor of https://github.com/cyberark/summon-aws-secrets/blob/master/main.go
// This was needed to have the command return the byte stream rather than have it write
// to STDOUT
func RetrieveSecret(variableName string) (secretBytes []byte, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)

	// Check if key has been specified
	arguments := strings.SplitN(variableName, "#", 2)

	secretName := arguments[0]
	var keyName string

	if len(arguments) > 1 {
		keyName = arguments[1]
	}

	exists, err := CheckIfSecretExists(secretName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("'%s' secret does not exist", secretName)
	}

	// Get secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
	  }

	resp, err := svc.GetSecretValue(context.TODO(), input)
    if err != nil {
		return nil, err
	}

	if resp.SecretString != nil {
		secretBytes = []byte(*resp.SecretString)
	} else {
		secretBytes = resp.SecretBinary
	}

	if keyName != "" {
		secretBytes, err = getValueByKey(keyName, secretBytes)
		if err != nil {
			return nil, err
		}
	}

	return
}

// ListSecrets will retrieval ALL secrets via pagination of 100 per page. It will
// return once all pages have been processed.
func ListSecrets() (secrets []types.SecretListEntry, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)

	// Get all secret names
    paginator := secretsmanager.NewListSecretsPaginator(svc, &secretsmanager.ListSecretsInput{
        MaxResults: aws.Int32(100),
    })

    for paginator.HasMorePages() {
        page, err := paginator.NextPage(context.TODO())
        if err != nil {
            return nil, fmt.Errorf("failed to retrieve secrets: %w", err)
        }

        secrets = append(secrets, page.SecretList...)
    }

    return secrets, nil
}

// GetSecret will retrieve a specific secret by Name (id)
func GetSecret(id string) (secret *secretsmanager.GetSecretValueOutput, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)


	secret, err = svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return
}

// DeleteSecret will retrieve a specific secret by Name (id)
func DeleteSecret(id string, force bool) (secret *secretsmanager.DeleteSecretOutput, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)

	secret, err = svc.DeleteSecret(context.TODO(), &secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(id),
		ForceDeleteWithoutRecovery: aws.Bool(force),
	})

	if err != nil {
		return nil, err
	}

	return
}

// PutSecretString will put an updated SecretString value to a specific secret by Name (id)
func PutSecretString(id string, data string) (secret *secretsmanager.PutSecretValueOutput, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)


	secret, err = svc.PutSecretValue(context.TODO(), &secretsmanager.PutSecretValueInput{
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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)

	secret, err = svc.PutSecretValue(context.TODO(), &secretsmanager.PutSecretValueInput{
		SecretBinary: data,
		SecretId:     aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return
}

// CreateSecretString will create a new SecretString value to a specific secret by Name (id)
func CreateSecretString(id string, data string, description string, tagsCSV string) (secret *secretsmanager.CreateSecretOutput, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)


	input := secretsmanager.CreateSecretInput{
		SecretString: aws.String(data),
		Name:         aws.String(id),
	}

	if description != "" {
		input.Description = aws.String(description)
	}

	if tagsCSV != "" {
		tags := []types.Tag{}
		for _, kv := range strings.Split(tagsCSV, ",") {
			parts := strings.SplitN(kv, "=", 2)
			tags = append(tags, types.Tag{
				Key:   aws.String(parts[0]),
				Value: aws.String(parts[1]),
			})
		}
		input.Tags = tags
	}

	secret, err = svc.CreateSecret(context.TODO(), &input)

	if err != nil {
		return nil, err
	}

	return
}

// CreateSecretBinary will create a new SecretBinary value to a specific secret by Name (id)
func CreateSecretBinary(id string, data []byte, description string, tagsCSV string) (secret *secretsmanager.CreateSecretOutput, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)


	input := secretsmanager.CreateSecretInput{
		SecretBinary: data,
		Name:         aws.String(id),
	}

	if description != "" {
		input.Description = aws.String(description)
	}

	if tagsCSV != "" {
		tags := []types.Tag{}
		for _, kv := range strings.Split(tagsCSV, ",") {
			parts := strings.SplitN(kv, "=", 2)
			tags = append(tags, types.Tag{
				Key:   aws.String(parts[0]),
				Value: aws.String(parts[1]),
			})
		}
		input.Tags = tags
	}

	secret, err = svc.CreateSecret(context.TODO(), &input)

	if err != nil {
		return nil, err
	}

	return
}

// DescribeSecret retrieves the describe data for a specific secret by Name (id)
func DescribeSecret(id string) (secret *secretsmanager.DescribeSecretOutput, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(cfg)

	secret, err = svc.DescribeSecret(context.TODO(), &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		return nil, err
	}

	return
}

// CheckIfSecretExists determines if the input secret ID already exists in AWS Secrets Manager
func CheckIfSecretExists(id string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return true, err
	}
	svc := secretsmanager.NewFromConfig(cfg)

	_, err = svc.DescribeSecret(context.TODO(), &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(id),
	})

	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return false, nil
		}
		return true, err
	}

	return true, nil
}
