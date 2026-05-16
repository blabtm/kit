#ifndef V2K_EM_ES_H
#define V2K_EM_ES_H

#include <ceres/ceres.h>
#include <v2k/ccd/config.h>
#include <v2k/em_es/config.h>
#include <yaml-cpp/yaml.h>

namespace v2k {
namespace em_es {

struct AxisRealization {
  v2k::ccd::X const &config;
  v2k::em_es::X const &setup;

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

}; // namespace em_es
}; // namespace v2k

#endif // V2K_EM_ES_H
