cmake_minimum_required(VERSION 3.11)
project(v2k CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS true)
set(CMAKE_BUILD_TYPE Release)

find_package(spdlog CONFIG REQUIRED)
find_package(GTest CONFIG REQUIRED)

enable_testing()
include(GoogleTest)

add_subdirectory(extern)
add_subdirectory(ccd)
add_subdirectory(em_es)

