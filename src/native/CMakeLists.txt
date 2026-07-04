cmake_minimum_required(VERSION 3.28)
project(lab C CXX)

add_compile_options(-Wno-error=date-time)

set(CMAKE_C_STANDARD 23)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(CMAKE_EXPORT_COMPILE_COMMANDS true)
set(CMAKE_BUILD_TYPE Release)

set(CMAKE_POLICY_VERSION_MINIMUM 3.10)
add_subdirectory(extern)

add_subdirectory(bootstrap_c)
