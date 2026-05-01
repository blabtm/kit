#ifndef V2K_EM_ES_H
#define V2K_EM_ES_H

#include <ceres/ceres.h>
#include <v2k/ccd.h>
#include <yaml-cpp/yaml.h>

namespace v2k {
namespace ems {

struct AxisSetup {
  double weight;

  AxisSetup(YAML::Node const &node) : weight(node["weight"].as<double>()) {}
};

struct AxesSetup {
  AxisSetup x;
  AxisSetup z;

  AxesSetup(YAML::Node const &node)
      : x(node["x"]), z(node["z"]) {}
};

struct CamSetup {
  AxesSetup axes;

  CamSetup(YAML::Node const &node) : axes(node["axes"]) {}
};

struct Config {
  std::size_t window;
  std::unordered_map<std::string, CamSetup> cams;

  /**
   *  Initialize configuration from specified location.
   */
  Config(std::string const &path) {
    auto const conf = YAML::LoadFile(path);
    auto const ccds = conf["service"]["cams"];

    window = conf["service"]["window"].as<std::size_t>();

    for (YAML::const_iterator it = ccds.begin(); it != ccds.end(); ++it) {
      cams.insert({it->first.as<std::string>(), CamSetup(it->second)});
    }
  }
};

struct AxisRealization {
  v2k::ccd::AxisConfig const &config;
  v2k::ems::AxisSetup const &setup;

  double measurement;
};

class CostFunction : public ceres::CostFunction {
private:
  std::vector<AxisRealization> const &axes_;

public:
  CostFunction(std::vector<AxisRealization> const &axes) : axes_(axes) {
    set_num_residuals(axes.size());
    mutable_parameter_block_sizes()->push_back(2);
  };

  ~CostFunction() {}

  bool Evaluate(double const *const *x, double *r, double **j) const override;
};

bool Estimate(std::vector<AxisRealization> const &axes, double *x);

}; // namespace ems
}; // namespace v2k

#endif // V2K_EM_ES_H
