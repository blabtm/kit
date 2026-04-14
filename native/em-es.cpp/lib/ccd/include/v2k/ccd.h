#ifndef V2K_CCD_H
#define V2K_CCD_H

#include <string>
#include <unordered_map>
#include <yaml-cpp/yaml.h>

namespace v2k {
namespace ccd {

struct AxisConfig {
  double beta;
  double dispersion;

  AxisConfig(YAML::Node const &node)
      : beta(node["beta"].as<double>()),
        dispersion(node["dispersion"].as<double>()) {}
};

struct AxesConfig {
  AxisConfig x;
  AxisConfig z;

  AxesConfig(YAML::Node const &node)
      : x(node["x"]), z(node["z"]) {}
};

struct CamConfig {
  AxesConfig axes;

  CamConfig(YAML::Node const &node) : axes(node["axes"]) {}
};

struct Config {
  std::unordered_map<std::string, CamConfig> cams;

  /**
   *  Initialize configuration from specified location.
   */
  Config(std::string const &path) {
    auto const conf = YAML::LoadFile(path);
    auto const ccds = conf["cams"];

    for (YAML::const_iterator it = ccds.begin(); it != ccds.end(); ++it) {
      cams.insert({it->first.as<std::string>(), CamConfig(it->second)});
    }
  }
};

}; // namespace ccd
}; // namespace v2k

#endif // V2K_CCD_H
