---
sidebar_position: 1
---

# Overview

## Administration

### Services

A service is a standalone deployment unit defined on the platform. The simplest service definition looks like this:

```yaml
driver: docker
```

Here we specify that the service is a Docker container.

#### Driver

A _service driver_ is an isolated deployment environment. The driver abstraction includes the following methods:

- **Ps** -- get current service status
- **Up** -- start the service
- **Down** -- stop the service

Currently, the platform supports the following set of drivers:

- `docker` -- containerized application in a Swarm cluster

Driver selection and preparation of all necessary artifacts (e.g., Dockerfile creation) is the developer's responsibility. From the administrator's perspective, all services are hidden behind the driver abstraction.

#### Service Manager

The _service manager_ is middleware that provides unified access to platform services and their configurations using the driver abstraction.

The service manager provides a REST API and a command-line interface. See the [service manager documentation](@site/docs/infra/sm/overview.md) for details.

### Build and Deployment

#### Build System

The platform repository is polyglot. Each language family uses its own specific build system. However, the administrator is not required to know all build systems in the project. To achieve this goal, a simple build layer was introduced: [Taskfile](https://taskfile.dev).

To view all available tasks and their descriptions:

```bash
task --list-all
```

For example, to deploy the platform, run:

```bash
task deploy
```

In addition to utility tasks, each service defines its own tasks for building and launching. The task structure fully mirrors the repository structure. Thus, for a service located in `jvm/em-es`, tasks will be named `jvm:em-es:*`.

Most services provide an `image` command to build the corresponding Docker image. For example, to build the service manager image:

```bash
task platform:sm:image
```

#### Configuration

Service-specific configurations are located in the `config` directory. After deployment, this directory will be accessible over the network via NFS and mounted to all running services on the platform in read-only mode. This allows each service to access all its configurations through the filesystem.

Configuration changes made through the operator platform are made in the `config` directory. This provides a single source of truth and enables versioning (since the `config` directory is under Git version control). This allows capturing changes over a period of time, rolling back to a previous version, and synchronizing with a remote repository.

Service configurations should not contain sensitive information. The process for injecting such information into services is described below.

Platform service configurations are located in the `platform/config` directory. Unlike regular service configurations, which are mounted directly to each service, platform configurations are registered in Docker Swarm Secrets and made available to each service individually.

#### Variables and Secrets

For more flexible and secure infrastructure management, variables and secrets are provided. They are defined as environment variables and located in the `platform/env` directory:

- `secret.env` -- contains sensitive variables, such as database access credentials
- `build.env` -- contains regular variables

The difference between secrets and regular variables is that secrets are not indexed in VCS.

Variables and secrets are _injected_ into platform configuration and infrastructure files by running:

```bash
task sync
```

As noted above, service configuration files should not contain environment variables. However, if a service needs to receive a variable value, this can be done through infrastructure files. For example, for Docker services, environment variables will be injected into the corresponding `compose.yaml` files.

#### Deployment

To deploy the platform, run a single command:

```bash
task deploy
```

Required files will be synchronized, variables injected, and platform services started on the Docker Swarm cluster according to `platform/compose.yaml`.

## Operation Guide

Everything the VEPP-2000 collider operator needs is available directly from the operator portal. Each service is represented as a separate application, including its configuration, visualization, and documentation.
