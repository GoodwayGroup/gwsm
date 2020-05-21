package info

var ViewCommandHelpText = `
View the current environment variables for a given ConfigMap and summon secrets.yml.

This will retrieve the stored secrets within AWS Secrets Manager and map them via
the secrets.yml file used by the 'summon' CLI tool to generate the current state of
Environment Variables for a given stage.

The AWS Secrets Manager names are assumed to be stored as '<SECRETS_GROUP>_NAME' 
in the ConfigMap. Example: 'RDS_SECRET_NAME: rds/staging/service-yolo'

From the root of the service, the required files are typically found below:

The path to the configmap.yaml file is within the kubernetes deployment.
This is typically .kube/<stage>/05-configmap.yaml

The path to the secrets.yml is typically .docker/secrets.yaml
`