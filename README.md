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
   1.1.0

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

### [`asdf` plugin](https://github.com/GoodwayGroup/asdf-gwsm)

Add plugin:

```
$ asdf plugin-add gwsm https://github.com/GoodwayGroup/asdf-gwsm.git
```

Install the latest version:

```
$ asdf install gwsm latest
```

### [Homebrew](https://brew.sh) (for macOS users)

```
brew tap GoodwayGroup/gwsm
brew install gwsm
```

### curl binary

```
$ curl https://i.jpillora.com/GoodwayGroup/gwsm! | bash
```

### [docker](https://www.docker.com/)

The compiled docker images are maintained on [GitHub Container Registry (ghcr.io)](https://github.com/orgs/GoodwayGroup/packages/container/package/gwsm).
We maintain the following tags:

- `edge`: Image that is build from the current `HEAD` of the main line branch.
- `latest`: Image that is built from the [latest released version](https://github.com/GoodwayGroup/gwsm/releases)
- `x.y.z` (versions): Images that are build from the tagged versions within Github.

```bash
docker pull ghcr.io/goodwaygroup/gwsm
docker run -v "$PWD":/workdir ghcr.io/goodwaygroup/gwsm --version
```

### man page

To install `man` page:

```
$ gwsm install-manpage
```

## Commands

Please see [the docs for details on the commands.](./docs/gwsm.md)

Each command ans subcommand has detailed help text that can be viewed using the `--help, -h` flag.

## Built With

- go v1.22
- make
- [git-chglog](https://github.com/git-chglog/git-chglog)
- [goreleaser](https://goreleaser.com/install/)

## Deployment

Run `./release.sh $VERSION`

This will update docs, changelog, add the tag, push main and the tag to the repo. The `goreleaser` action will publish the binaries to the Github Release.

If you want to simulate the `goreleaser` process, run the following command:

```
$ curl -sL https://git.io/goreleaser | bash -s -- --rm-dist --skip-publish --snapshot
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull
requests to us.

1. Fork the [GoodwayGroup/gwsm](https://github.com/GoodwayGroup/gwsm) repo
1. Use `go >= 1.16`
1. Branch & Code
1. Run linters :broom: `golangci-lint run`
   - The project uses [golangci-lint](https://golangci-lint.run/welcome/install/#brew)
1. Commit with a Conventional Commit
1. Open a PR

## Versioning

We employ [git-chglog](https://github.com/git-chglog/git-chglog) to manage the [CHANGELOG.md](CHANGELOG.md). For the
versions available, see the [tags on this repository](https://github.com/GoodwayGroup/gwsm/tags).

## Authors

- **Derek Smith** - [@clok](https://github.com/clok)

See also the list of [contributors](https://github.com/GoodwayGroup/gwvault/contributors) who participated in this
project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Sponsors

[![goodwaygroup][goodwaygroup]](https://goodwaygroup.com)

[goodwaygroup]: https://s3.amazonaws.com/gw-crs-assets/goodwaygroup/logos/ggLogo_sm.png "Goodway Group"
