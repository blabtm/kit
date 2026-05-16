find_package(yaml-cpp CONFIG REQUIRED)
find_package(nlohmann_json CONFIG REQUIRED)

add_library(ccd INTERFACE)
target_include_directories(ccd INTERFACE include)
target_link_libraries(ccd INTERFACE yaml-cpp::yaml-cpp)
target_link_libraries(ccd INTERFACE nlohmann_json::nlohmann_json)

