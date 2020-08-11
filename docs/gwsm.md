% gwsm 8
# NAME
gwsm - interact with config map and secret manager variables
# SYNOPSIS
gwsm

**Usage**:
```
gwsm [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## secretsmanager, sm

Secrets Manager commands w/ interactive interface

**--binary, -b**: get the SecretBinary value

### list

display table of all secrets with meta data

### describe

print description of secret to STDOUT

**--secret-id, -s**="": Specific Secret to describe, will bypass select/search

### get, view

select from list or pass in specific secret

**--secret-id, -s**="": Specific Secret to view, will bypass select/search

### edit, e

interactive edit of a secret String Value

**--secret-id, -s**="": Specific Secret to edit, will bypass select/search

### create, c

create new secret in Secrets Manager

**--description, --desc**="": Additional description text.

**--interactive, -i**: Open interactive editor to create secret value.

**--secret-id, -s**="": Secret name

**--tags**="": key=value tags (CSV list)

**--value, -v**="": Secret Value. Will store as a string, unless binary flag is set.

### put

non-interactive update to a specific secret

**--description, --desc**="": Additional description text.

**--interactive, -i**: Open interactive editor to create secret value.

**--secret-id, -s**="": Secret name

**--tags**="": key=value tags (CSV list)

**--value, -v**="": Secret Value. Will store as a string, unless binary flag is set.

### delete, del

delete a specific secret

**--force, -f**: Bypass recovery window (30 days) and immediately delete Secret.

**--secret-id, -s**="": Specific Secret to delete

## diff, d

View diff of local vs. namespace

```
View the diff of the local environment against a given command running on a
pod within a namespace.

This will retrieve the stored secrets within AWS Secrets Manager and map them
via the secrets.yml file used by the 'summon' CLI tool to generate the current
state of Environment Variables for a given stage.

The AWS Secrets Manager names are assumed to be stored as
'<SECRETS_GROUP>_NAME' in the ConfigMap.
Example: 'RDS_SECRET_NAME: rds/staging/service-yolo'

From the root of the service, the required files are typically found below:

The path to the configmap.yaml file is within the kubernetes deployment.
This is typically .kube/<stage>/05-configmap.yaml

The path to the secrets.yml is typically .docker/secrets.yaml

It will then grab current environment for a specific process running within a
Pod in a given Namespace.

This is achieved by inspecting the /proc/<PID>/environ for the given process.
This method uses '/bin/bash -c' as the base command to perform the PID
inspection via 'ps faux'.

The 'filter-prefix' flag will exclude any values that start with the flagged
prefixes from display.

The 'exclude' flag will exclude any values where the KEY matches exactly from
display.
```

**--cmd**="": Command to inspect (default: node)

**--configmap, -c**="": Path to configmap.yaml

**--exclude**="": List (csv) of specific env vars to exclude values from display. Set to "" to remove any exclusions. (default: PATH,SHLVL,HOSTNAME)

**--filter-prefix, -f**="": List of prefixes (csv) used to filter values from display. Set to "" to remove any filters. (default: npm_,KUBERNETES_,API_PORT)

**--namespace, -n**="": Kube Namespace list Pods from

**--secret-suffix**="": Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup (default: _NAME)

**--secrets, -s**="": Path to secrets.yml (default: .docker/secrets.yml)

## diff:legacy, diff:ansible

View diff of local (ansible encrypted) vs. namespace

```
View the diff of the local ansible-vault encrypted Kubenetes Secret file
against a given dotenv file on a pod within a namespace.

The local file will use the contents of the 'data.<accsessor flag>' block.
This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>

It will then grab contents of the dotenv filr on a Pod in a given Namespace.

This defaults to inspecting the '$PWD/.env on' when executing a 'cat' command.
This method uses '/bin/bash -c' as the base command to perform inspection.
```

**--accessor, -a**="": Accessor key to pull data out of Data block. (default: .env)

**--dotenv**="": Path to .env file on Pod (default: $PWD/.env)

**--encrypted-env-file, -e**="": Path to encrypted Kube Secret file

**--namespace, -n**="": Kube Namespace list Pods from

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## local, l

Interact with local env files

### view, v

View values based on local settings

```
View the current environment variables for a given ConfigMap and summon
secrets.yml.

This will retrieve the stored secrets within AWS Secrets Manager and map them
via the secrets.yml file used by the 'summon' CLI tool to generate the current
state of Environment Variables for a given stage.

The AWS Secrets Manager names are assumed to be stored as
'<SECRETS_GROUP>_NAME' in the ConfigMap. 
Example: 'RDS_SECRET_NAME: rds/staging/service-yolo'

From the root of the service, the required files are typically found below:

The path to the configmap.yaml file is within the kubernetes deployment.
This is typically .kube/<stage>/05-configmap.yaml

The path to the secrets.yml is typically .docker/secrets.yaml

The 'filter-prefix' flag will exclude any values that start with the flagged 
prefixes from display.

The 'exclude' flag will exclude any values where the KEY matches exactly from
display.
```

**--configmap, -c**="": Path to configmap.yaml

**--secret-suffix**="": Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup (default: _NAME)

**--secrets, -s**="": Path to secrets.yml (default: .docker/secrets.yml)

### ansible, legacy, a

View value from ansible-vault encrypted Kube Secret file.

```
View a legacy ansible-vault encrypted Kubenetes Secret file. This will output
the contents of the 'data.<accsessor flag>' block.
This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>
```

**--accessor, -a**="": Accessor key to pull data out of Data block. (default: .env)

**--encrypted-env-file, -e**="": Path to encrypted Kube Secret file

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## namespace, ns

Interact with env on a running Pod within a Namespace

### view, v

View values configured within a namespace

```
View the current environment for a specific process running within a Pod in a
given Namespace.

This is achieved by inspecting the /proc/<PID>/environ for the given process.
This method uses '/bin/bash -c' as the base command to perform the PID
inspection via 'ps faux'.

The 'filter-prefix' flag will exclude any values that start with the flagged
prefixes from display.

The 'exclude' flag will exclude any values where the KEY matches exactly from
display.
```

**--cmd**="": Command to inspect (default: node)

**--exclude**="": List (csv) of specific env vars to exclude values from display. Set to "" to remove any exclusions. (default: PATH,SHLVL,HOSTNAME)

**--filter-prefix, -f**="": List of prefixes (csv) used to filter values from display. Set to "" to remove any filters. (default: npm_,KUBERNETES_,API_PORT)

**--namespace, -n**="": Kube Namespace list Pods from

## s3

simple S3 commands

### get

[object path] [destination path]

```
The '[object path]' MUST always start with 's3://'
The '[destination path]' directory MUST exists, but file will be created or overwritten

Example:
$ gwsm s3 get s3://coll-bucket-name/with/path/filename /tmp/filename
```

## version, v

Print version info

