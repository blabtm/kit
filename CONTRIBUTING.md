# Contribution Guide

Here you can find an overall rules and practices to follow as the platform contributor. For details
on specific language development environment see the corresponding directory.

## Development Environment

### Requirements

- [Docker Engine](https://docs.docker.com/engine/install/)

### VS Code

To build and run container inside the VS Code you need to install an appropriate
[extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).

Start VS Code, run **Dev Containers: Open Folder in Container** from the Command Palette (`F1`) and
select the platform repository. After that, wait for container to be built (it may take a while, but
only for the first time).

For more information, see the [docs](https://code.visualstudio.com/docs/devcontainers/containers).

### CLI

If you prefer something else then VS Code or just want to run container using shell you can use the
[development container CLI](https://github.com/devcontainers/cli).

After installation complete, navigate to the project directory, start the container and attach to
it:

```bash
devcontainer up
devcontainer exec /bin/bash
```

## Project Structure

- _platform_ - platform middleware and infrastructure
- _config_ - service configuration files
- _tools_ - build utilities and extensions
- _scripts_ - helper scripts
- _schema_ - protocol definitions
- _native_ - C/C++ services
- _py_ - Python services
- _go_ - Golang services
- _jvm_ - Java, Scala, Kotlin and Groovy services
