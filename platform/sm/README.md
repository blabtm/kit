# Service Manager

_Service Manager_ is a core middleware of
[the VEPP-2000 Collider Platform](https://github.com/blabtm/v2k). It provides unified way to control
and monitor heterogeneous services defined on the platform.

## REST API

The primary access point is the REST API. If you don't know what the REST API is, see
[the Red Hat explanation](https://www.redhat.com/en/topics/api/what-is-a-rest-api).

## CLI

The CLI is a convenient way to access Service Manager from the shell. The recommended way to install
it is:

```bash
go install github.com/blabtm/v2k/platform/sm@v0.0.0
```

By default, CLI will act as a REST client to the existing Service Manager instance. In case you
don't have one running, you can use CLI in standalone mode.

You can see all available commands, their description and flags with:

```bash
sm help
```
