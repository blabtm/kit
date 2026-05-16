find_package(protobuf CONFIG REQUIRED)
find_package(Ceres CONFIG REQUIRED)
find_package(unofficial-modern-cpp-kafka CONFIG REQUIRED)

add_library(em-es
  src/em_es.cc
)

target_include_directories(em-es PUBLIC include)
target_link_libraries(em-es PUBLIC spdlog::spdlog)
target_link_libraries(em-es PUBLIC Ceres::ceres)
target_link_libraries(em-es PUBLIC ccd)

add_executable(svc-em-es
  src/main.cc
  gen/v2k/Log.pb.cc
  gen/v2k/beam/AxisSize.pb.cc
  gen/v2k/beam/Current.pb.cc
  gen/v2k/beam/Emittance.pb.cc
  gen/v2k/beam/EnergySpread.pb.cc
)

target_include_directories(svc-em-es PRIVATE gen)
target_link_libraries(svc-em-es PUBLIC spdlog::spdlog)
target_link_libraries(svc-em-es PRIVATE protobuf::libprotobuf)
target_link_libraries(svc-em-es PRIVATE unofficial::modern-cpp-kafka::modern-cpp-kafka)
target_link_libraries(svc-em-es PRIVATE em-es)
target_link_libraries(svc-em-es PRIVATE ccd)
