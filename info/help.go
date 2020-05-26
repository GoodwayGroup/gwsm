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
