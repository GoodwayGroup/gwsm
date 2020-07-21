<a name="unreleased"></a>
## [Unreleased]


<a name="v0.8.1"></a>
## [v0.8.1] - 2020-07-21
### Chore
- update dependencies

### DevOps
- updated release.sh script

### Features
- **release:** v0.8.0


<a name="v0.8.0"></a>
## [v0.8.0] - 2020-07-20
### Features
- **awssession:** use awssession to manage session creation
- **release:** v0.7.2


<a name="v0.7.2"></a>
## [v0.7.2] - 2020-06-10
### Bug Fixes
- **sm:** fixed hanging goroutine when error returned from SM

### Features
- **release:** v0.7.1


<a name="v0.7.1"></a>
## [v0.7.1] - 2020-06-09
### Debt
- **version:** use build flag instead of static file for versioning

### Features
- **release:** v0.7.0


<a name="v0.7.0"></a>
## [v0.7.0] - 2020-06-09
### Features
- **release:** v0.6.1
- **sm:** support for delete secret, deprecated custom editor code in favor of survey prompt


<a name="v0.6.1"></a>
## [v0.6.1] - 2020-06-08
### Chore
- **version:** bump version of app to v0.6.0

### Features
- **release:** v0.6.0
- **sm:** added JSON validation step to edit action that allows for verification before submit


<a name="v0.6.0"></a>
## [v0.6.0] - 2020-06-08
### Debt
- **logger:** updated logging and added tests

### Features
- **release:** v0.5.0
- **sm:** Added interactive editor for edit and create
- **sm:** Added list, describe, get, edit and create SecretsManager commmands


<a name="v0.5.0"></a>
## [v0.5.0] - 2020-06-03
### Bug Fixes
- **typo:** address typo in error message

### Chore
- **changelog:** updated settings for changelog generation
- **version:** bump app version to v0.5.0

### Features
- **release:** v0.4.0
- **s3:** added simple s3 get command


<a name="v0.4.0"></a>
## [v0.4.0] - 2020-05-27
### Chore
- **verion:** bump version to v0.4.0

### Features
- **release:** v0.3.0

### Pull Requests
- feat(ansible): support for ansible-vault encrypted Kube Secret files ([#2](https://github.com/GoodwayGroup/gwsm/issues/2))


<a name="v0.3.0"></a>
## [v0.3.0] - 2020-05-26
### Chore
- **version:** bump version to v0.3.0

### Pull Requests
- feat(diff): View local and environment running on Pod with support for diff between the two. ([#1](https://github.com/GoodwayGroup/gwsm/issues/1))


<a name="v0.2.1"></a>
## [v0.2.1] - 2020-05-21
### Chore
- **app:** adjust app name :facepalm:
- **name:** adjust name in all the right places

### Features
- **release:** v0.2.0


<a name="v0.2.0"></a>
## [v0.2.0] - 2020-05-21
### Chore
- **version:** bump version to v0.2.0

### Features
- **pod:** add new kube command files
- **pod:** view environment for a command running within a pod
- **release:** v0.1.0


<a name="v0.1.0"></a>
## v0.1.0 - 2020-05-21
### Bug Fixes
- **release:** swapped tagging and CHANGELOG generation

### Features
- **view:** Added view command and scaffolding for deployments


[Unreleased]: https://github.com/GoodwayGroup/gwsm/compare/v0.8.1...HEAD
[v0.8.1]: https://github.com/GoodwayGroup/gwsm/compare/v0.8.0...v0.8.1
[v0.8.0]: https://github.com/GoodwayGroup/gwsm/compare/v0.7.2...v0.8.0
[v0.7.2]: https://github.com/GoodwayGroup/gwsm/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/GoodwayGroup/gwsm/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/GoodwayGroup/gwsm/compare/v0.6.1...v0.7.0
[v0.6.1]: https://github.com/GoodwayGroup/gwsm/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/GoodwayGroup/gwsm/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/GoodwayGroup/gwsm/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/GoodwayGroup/gwsm/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/GoodwayGroup/gwsm/compare/v0.2.1...v0.3.0
[v0.2.1]: https://github.com/GoodwayGroup/gwsm/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/GoodwayGroup/gwsm/compare/v0.1.0...v0.2.0
