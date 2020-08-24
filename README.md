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

VERSION:
   v0.12.4

AUTHORS:
   Derek Smith <dsmith@goodwaygroup.com>
   Goodway Group Tech

COMMANDS:
   sm, secretsmanager  Secrets Manager commands w/ interactive interface
   env, e              Commands to interact with environment variables, both local and on cluster.
   s3                  simple S3 commands
   install-manpage     Generate and install man page
   version, v          Print version info
   help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

COPYRIGHT:
   (c) 2020 Goodway Group
```

## Installation

To install latest version:

```
$ curl https://i.jpillora.com/GoodwayGroup/gwsm! | bash
```

To install `man` page:

```
$ gwsm install-manpage
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
