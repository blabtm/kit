#include <spdlog/spdlog.h>
#include <v2k/em_es/em_es.h>

bool v2k::em_es::CostFunction::Evaluate(double const *const *x, double *r,
                                        double **j) const {
  for (int i = 0; i < axes_.size(); ++i) {
    auto const &axis = axes_[i];

    double const size = axis.measurement;
    double const beta = axis.config.get_beta();
    double const dispersion = axis.config.get_dispersion();
    double const emittance = x[0][0];
    double const spread = x[0][1];
    double const dxs = dispersion * spread;

    r[i] = (size * size) - (beta * emittance) - (dxs * dxs);

    if (j != nullptr && j[0] != nullptr) {
      j[0][i * 2 + 0] = -beta;
      j[0][i * 2 + 1] = -2 * dxs * dispersion;
    }
  }

  return true;
};

bool v2k::em_es::Estimate(std::vector<AxisRealization> const &axes, double *x) {
  v2k::em_es::CostFunction *const function = new CostFunction(axes);
  ceres::Problem problem;

  problem.AddResidualBlock(function, nullptr, x);
  problem.SetParameterLowerBound(x, 0, 0);
  problem.SetParameterLowerBound(x, 1, 0);

  ceres::Solver::Options options;
  options.linear_solver_type = ceres::DENSE_QR;
  ceres::Solver::Summary summary;
  ceres::Solve(options, &problem, &summary);

  spdlog::debug(summary.BriefReport());
  spdlog::debug("Emittance: {:.6e}, Energy Spread: {:.6e}", x[0], x[1]);

  return summary.IsSolutionUsable();
}
