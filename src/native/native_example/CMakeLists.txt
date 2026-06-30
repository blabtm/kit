add_executable(native_example
  src/main.c
)

target_link_libraries(native_example PUBLIC zlog)
