# Timber.io - Master your Go apps with structured logging

Timber Go is under development, currently this library only contains the core components to ship logs to the Timber API.

## Packages

### `batch`

An implementation for efficiently collecting and sending strings via input and output channels. Batches are sent when the configured buffer size is exceeded or if the configured time has elasped since the last send.

### `forward`

Exposes an a Forwarder interface for accepting a buffer and writing it somewhere. Implementations include stdout, file,
and http forwarders.

### `logging`

Exposes a logger interface that all timber-go packages and their dependencies use. Aside from bring your own, there is
the default logger to stdout and the discard logger to /dev/null. This allows for deciding how library logs are handled
within their dependents.

### `metadata`

An implementation of Timber's metadata JSON schema. It also includes metadata collection for the supported platforms
(e.g. AWS).

## LICENSE

The original parts of this software as developed by Timber Technologies, Inc. as
well as contributors are licensed under the Internet Systems Consortium (ISC)
License. This software is dependent upon third-party code which is
statically linked into the executable at compile time. This third-party code is
redistributed without modification and made available to all users  under the
terms of each project's original license within the `vendor` directory of the
project.
