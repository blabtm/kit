add_subdirectory(zlog)
target_include_directories(zlog PUBLIC zlog/src)
target_include_directories(zlog_s PUBLIC zlog/src)

add_subdirectory(gtest)

set(protobuf_BUILD_TESTS OFF CACHE BOOL "" FORCE)
set(protobuf_BUILD_SHARED_LIBS OFF CACHE BOOL "" FORCE)
add_subdirectory(protobuf)

set(Protobuf_FOUND TRUE CACHE INTERNAL "")
set(Protobuf_INCLUDE_DIRS
  "${CMAKE_CURRENT_SOURCE_DIR}/protobuf/src"
  CACHE INTERNAL "")
set(Protobuf_INCLUDE_DIR 
  "${CMAKE_CURRENT_SOURCE_DIR}/protobuf/src"
  CACHE INTERNAL "")
set(Protobuf_LIBRARIES
  protobuf::libprotobuf
  CACHE INTERNAL "")
set(Protobuf_LIBRARY
  protobuf::libprotobuf
  CACHE INTERNAL "")
set(Protobuf_PROTOC_EXECUTABLE
  $<TARGET_FILE:protobuf::protoc>
  CACHE INTERNAL "")

set(Protobuf_VERSION "7.35.1" CACHE INTERNAL "")
set(Protobuf_VERSION_MAJOR "7" CACHE INTERNAL "")
set(Protobuf_VERSION_MINOR "35" CACHE INTERNAL "")
set(Protobuf_VERSION_PATCH "1" CACHE INTERNAL "")

set(CMAKE_SKIP_INSTALL_RULES ON)

add_subdirectory(protobuf-c/build-cmake)
