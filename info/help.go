package info

var ViewCommandHelpText = `
View the Environment for
- local:     Based on the specified ConfigMap and summon secrets.yml
- namespace: Inspect the environment for a specific process running on a Pod
- diff:      Compare 'namespace' environment with 'local'
`

var ViewLocalCommandHelpText = `
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

The 'filter-prefix' flag will exclude any values that start with the flagged prefixes from display.

The 'exclude' flag will exclude any values where the KEY matches exactly from display.
`

var ViewNamespaceCommandHelpText = `
View the current environment for a specific process running within a Pod in a given Namespace.

This is achieved by inspecting the /proc/<PID>/environ for the given process. This method uses
'/bin/bash -c' as the base command to perform the PID inspection via 'ps faux'.

The 'filter-prefix' flag will exclude any values that start with the flagged prefixes from display.

The 'exclude' flag will exclude any values where the KEY matches exactly from display.
`

var ViewDiffCommandHelpText = `
View the diff of the local environment against a given command running on a pod within a namespace.

This will retrieve the stored secrets within AWS Secrets Manager and map them via
the secrets.yml file used by the 'summon' CLI tool to generate the current state of
Environment Variables for a given stage.

The AWS Secrets Manager names are assumed to be stored as '<SECRETS_GROUP>_NAME' 
in the ConfigMap. Example: 'RDS_SECRET_NAME: rds/staging/service-yolo'

From the root of the service, the required files are typically found below:

The path to the configmap.yaml file is within the kubernetes deployment.
This is typically .kube/<stage>/05-configmap.yaml

The path to the secrets.yml is typically .docker/secrets.yaml

It will then grab current environment for a specific process running within a Pod in a given Namespace.

This is achieved by inspecting the /proc/<PID>/environ for the given process. This method uses
'/bin/bash -c' as the base command to perform the PID inspection via 'ps faux'.

The 'filter-prefix' flag will exclude any values that start with the flagged prefixes from display.

The 'exclude' flag will exclude any values where the KEY matches exactly from display.
`

var ViewAnsibleEncryptedEnvCommandHelpText = `
View a legacy ansible-vault encrypted Kubenetes Secret file. This will output the contents of
the 'data.<accsessor flag>' block. This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>
`

var ViewAnsibleEnvDiffCommandHelpText = `
View the diff of the local ansible-vault encrypted Kubenetes Secret file against a given dotenv
file on a pod within a namespace.

The local file will use the contents of the 'data.<accsessor flag>' block. This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>

It will then grab contents of the dotenv filr on a Pod in a given Namespace.

This defaults to inspecting the '$PWD/.env on' when executing a 'cat' command. This method uses
'/bin/bash -c' as the base command to perform inspection.
`

var S3GetCommandHelp = `
The '[object path]' MUST always start with 's3://'
The '[destination path]' direcotry MUST exists, but file will be created or overwritten

Example:
$ gwsm s3 get s3://coll-bucket-name/with/path/filename /tmp/filename
`
