# GW Secrets Manager for Kubernetes

Kubernetes ConfigMaps + [summon](https://github.com/cyberark/summon) + AWS Secrets Manager

## Basic Usage

Please see [the docs for details on the commands.](./docs/gwsm.md)

```
$ gwsm help
NAME:
   gwsm - interact with config map and secret manager variables

USAGE:
   gwsm [global options] command [command options] [arguments...]

AUTHORS:
   Derek Smith <dsmith@goodwaygroup.com>
   Goodway Group Tech

COMMANDS:
   secretsmanager, sm         Secrets Manager commands w/ interactive interface
   diff, d                    View diff of local vs. namespace
   diff:legacy, diff:ansible  View diff of local (ansible encrypted) vs. namespace
   local, l                   Interact with local env files
   namespace, ns              Interact with env on a running Pod within a Namespace
   s3                         simple S3 commands
   version, v                 Print version info
   help, h                    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

COPYRIGHT:
   (c) 2020 Goodway Group


```

## Installation

```
$ curl https://i.jpillora.com/GoodwayGroup/gwsm! | bash
```

## Commands

Please see [the docs for details on the commands.](./docs/gwsm.md)

Each command ans subcommand has detailed help text that can be viewed using the `--help, -h` flag.


## Built With

* go v1.14+
* make
* [github.com/mitchellh/gox](https://github.com/mitchellh/gox)

## Deployment

Run `./release.sh $VERSION`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We employ [git-chglog](https://github.com/git-chglog/git-chglog) to manage the [CHANGELOG.md](CHANGELOG.md). For the versions available, see the [tags on this repository](https://github.com/GoodwayGroup/gwsm/tags).

## Authors

* **Derek Smith** - [@clok](https://github.com/clok)

See also the list of [contributors](https://github.com/GoodwayGroup/gwvault/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Sponsors

[![goodwaygroup][goodwaygroup]](https://goodwaygroup.com)

[goodwaygroup]: https://s3.amazonaws.com/gw-crs-assets/goodwaygroup/logos/ggLogo_sm.png "Goodway Group"
