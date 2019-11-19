# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

* Add a common interface to read partially hydrated workflow specifications.
* Add a client to retrieve secrets from the metadata API.

## [1.3.0]

### Added

* Add logging interface through go sdk and `ni log` command

## [1.2.0]

### Added

* Add support for validating Nebula workflow files.

## [1.1.1]

### Changed

* Fix: Clear the entrypoint in generated Docker images to align overridden
  Alpine Linux-based image behavior to the standard Alpine Linux image.

## [1.1.0]

### Added

* Initial implementation of the Spindle utility.

## [1.0.1]

### Changed

* Fix: Link `ni` binaries statically by setting `CGO_ENABLED=0`.

## [1.0.0]

### Added

* Public release of the `ni` command.

[Unreleased]: https://github.com/puppetlabs/nebula-sdk/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/puppetlabs/nebula-sdk/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/puppetlabs/nebula-sdk/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/puppetlabs/nebula-sdk/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/puppetlabs/nebula-sdk/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/puppetlabs/nebula-sdk/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/puppetlabs/nebula-sdk/compare/902be9735b850b21229bf34ddf42a11aba6b315e...v1.0.0
