# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2018-06-06

### Changed

  - Use `[]byte` instead of `string` for batcher input to avoid unnecessary overhead.
  - Discard httpClient logs as they were noisy and redundant.

## [0.0.1] - 2018-05-11

### Added

  - Imported core functionality for handling logs in memory from [Timber Agent](https://github.com/timberio/agent)

[Unreleased]: https://github.com/timberio/timber-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/timberio/timber-go/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/timberio/timber-go/tree/v0.0.1
