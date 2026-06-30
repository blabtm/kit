# Overview

The _Service Manager_ is the core middleware component of the
[VEPP-2000 collider platform](https://github.com/blabtm/v2k). It provides a unified way to manage
and monitor heterogeneous services defined on the platform.

## REST API

The primary entry point is the REST API. If you are unfamiliar with REST APIs, see the
[Red Hat explanation](https://www.redhat.com/en/topics/api/what-is-a-rest-api).

## CLI

The CLI provides a convenient way to access the Service Manager from the shell. The recommended installation method is:

```bash
go install github.com/blabtm/v2k/platform/sm@v0.0.0
```

By default, the CLI will act as a REST client to an existing Service Manager instance. If you
do not have an instance running, you can use the CLI in standalone mode.

You can view all available commands, their descriptions, and flags with:

```bash
sm help
```