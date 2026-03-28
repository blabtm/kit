# VEPP-2000 Collider Platform

## Administration Guide

### Service

Service is a single deployment unit defined on the platform. The simplest service definition is:

```yaml
driver: docker
```

Here we say, that the service is a docker container. That's all.

#### Driver

_Service Driver_ is an isolated deployment environment. Driver's abstraction includes the following
methods:

- **Ps** - get service's current status
- **Up** - trigger service to start
- **Down** - trigger service to stop

At the moment, platform supports the following set of drivers:

- `docker` - containerized application on the Swarm cluster
- `flink` - job on the Flink cluster (**in development**)

It's the developer's job to choose the driver and provide all the necessary artifacts (e.g. create
Dockerfile or assemble fat-jar). From the administrator's perspective, all the services are hidden
behind the driver's abstraction.

#### Service Manager

_Service Manager_ is a middleware to provide unified access to platform services and their
configurations using driver's abstraction.

Service Manager provides REST API and command line interface. For more information, see the
[Service Manager Docs](platform/sm/README.md)

### Build System

The platform repository is polyglot. Each language family employing it's own, specific build system.
Despite this, no knowledge about all the build systems in the project required from the
administrator. To achieve this, simple build layer was introduced: [Taskfile](https://taskfile.dev).

You can list all available tasks and there description:

```bash
task --list-all
```

For example, to inject the secrets and deploy the platform, one can run this:

```bash
task platform:deploy
```

### Configuration

Service-specific configurations are located in the `config` directory. Platform-specific
configurations are located in the `platform/config` directory. Platform secrets are located in the
`platform/secret/secret.env` file.

## Operation Guide

Everything you need as a VEPP-2000 collider operator is available directly from the operator's
platform. Each service is represented as a single app, including it's configurations, related
visualizations and documentation.
