# GW Secrets Mananger for Kubernetes

Kubernetes ConfigMaps + [summon](https://github.com/cyberark/summon) + AWS Secrets Manager

## Basic Usage

```
$ gwsm help
NAME:
   gwsm - interact with config map and secret manager variables

USAGE:
   gwsm [global options] command [command options] [arguments...]

VERSION:
   v0.1.0

AUTHOR:
   Derek Smith <dsmith@goodwaygroup.com>

COMMANDS:
   version, v  Print version info
   view        View the current environment variables for a given configmap and secrets.yml
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --region value  AWS region (default: us-east-1) [$AWS_REGION]
   --help, -h      show help (default: false)
   --version, -v   print the version (default: false)

COPYRIGHT:
   (c) 2020 Goodway Group
```

## Installation

```
$ curl https://i.jpillora.com/GoodwayGroup/gwsm! | bash
```

## Commands

### view

```
$ gwsm help view
NAME:
   gwsm view - View the current environment variables for a given configmap and secrets.yml

USAGE:

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


OPTIONS:
   --secrets value, -s value    Path to secrets.yml (default: ".docker/secrets.yml")
   --configmap value, -m value  Path to configmap.yaml
```


## Built With

* go v1.14+
* make
* [github.com/mitchellh/gox](https://github.com/mitchellh/gox)

## Deployment

Run `./release.sh $VERSION`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We employ [auto-changelog](https://www.npmjs.com/package/auto-changelog) to manage the [CHANGELOG.md](CHANGELOG.md). For the versions available, see the [tags on this repository](https://github.com/GoodwayGroup/gwvault/tags).

## Authors

* **Derek Smith** - [@clok](https://github.com/clok)

See also the list of [contributors](https://github.com/GoodwayGroup/gwvault/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Sponsors

[![goodwaygroup][goodwaygroup]](https://goodwaygroup.com)

[goodwaygroup]: https://s3.amazonaws.com/gw-crs-assets/goodwaygroup/logos/ggLogo_sm.png "Goodway Group"
