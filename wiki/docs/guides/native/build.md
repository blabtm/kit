---
sidebar_position: 1
---

# Система сборки

Основная система сборки: [CMake](https://cmake.org). Библиотеки и приложения, разрабатываемые в рамках проекта могут и должны переиспользоватья с целью минимизации дублирующегося функционала. 

Для получения и управления внешними зависимостями в проекте используется [vcpkg](https://vcpkg.io), который находится в репозитории в виде Git подмодуля по пути `native/extern/vcpkg`. Инструмент автоматически устанавливается при сборки контейнера разработки и доступен во всех проектах. Далее описан стандартный процесс работы.

1. Инициализация сервиса

    - Создание директории с соответствующим названием (см. [соглашения](../../conventions.md))

    ```bash
    mkdir native/em-es.cpp && cd native/em-es.cpp
    ```

    - Инициализация `vcpkg`:

    ```bash
    vcpkg new --application
    ```

2. Объявление внешних зависимостей

    Нужные зависимости можно найти в [регистре](https://vcpkg.io/en/packages).

    ```bash
    vcpkg add port protobuf
    ```

3. Инициализация CMake

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

4. Использование внутренних зависимостей

    Внутренние зависимости подключаются в виде CMake модулей:

    ```cmake
    add_subdirectory(../other)
    ```

5. Сборка

    ```bash
    mkdir build && cd build
    cmake .. && make
    ```