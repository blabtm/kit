---
sidebar_position: 1
---

# Build System

The primary build system is [CMake](https://cmake.org). Libraries and applications developed as part of the project can and should be reused to minimize duplicate functionality.

The project uses [vcpkg](https://vcpkg.io) for managing external dependencies, which is included in the repository as a Git submodule at `native/extern/vcpkg`. The tool is automatically installed when building the development container and is available in all projects. The standard workflow is described below.

1. Service Initialization

    - Create a directory with the appropriate name (see [conventions](../../conventions.md))

    ```bash
    mkdir native/em-es.cpp && cd native/em-es.cpp
    ```

    - Initialize `vcpkg`:

    ```bash
    vcpkg new --application
    ```

2. Declaring External Dependencies

    Required dependencies can be found in the [registry](https://vcpkg.io/en/packages).

    ```bash
    vcpkg add port protobuf
    ```

3. CMake Initialization

    ```cmake
    cmake_minimum_required(VERSION 3.11)
    set(CMAKE_TOOLCHAIN_FILE $ENV{VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake 
        CACHE STRING "Vcpkg toolchain file")
    project(em-es.cpp CXX)

    set(CMAKE_CXX_STANDARD 20)
    set(CMAKE_CXX_STANDARD_REQUIRED ON)
    set(CMAKE_EXPORT_COMPILE_COMMANDS true)
    set(CMAKE_BUILD_TYPE Debug)

    find_package(protobuf CONFIG REQUIRED)

    ...
    ```

4. Using Internal Dependencies

    Internal dependencies are included as CMake modules:

    ```cmake
    add_subdirectory(../other)
    ```

5. Building

    ```bash
    mkdir build && cd build
    cmake .. && make
    ```