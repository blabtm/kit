add_executable(bootstrap_c
  src/main.c
)

target_link_libraries(bootstrap_c PRIVATE
  zlog
  protobuf-c
)
