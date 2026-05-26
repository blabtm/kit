#include <absl/log/log.h>
#include <cmath>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>
#include <v2k/ccd.h>
#include <v2k/em_es.h>

namespace v2k {
namespace ems {

class EstimationTest : public testing::Test {
protected:
  EstimationTest() {}
  ~EstimationTest() override {}

  void SetUp() override {
    v2k::ccd::Config config("/workspaces/v2k/config/ccd/config.yaml");
  }

  void TearDown() override {}
};

TEST_F(EstimationTest, EstimateSyntheticData) {
  v2k::ccd::Config ccd_conf("/workspaces/v2k/config/ccd/config.yaml");
  v2k::ems::Config ems_conf("/workspaces/v2k/config/em-es/config.yaml");

  std::vector<AxisRealization> realizations;

  double const noise = 0.05;
  double const em = 4e-6;
  double const es = 6e-4;

  for (auto const &[name, setup] : ems_conf.cams) {
    auto const &conf = ccd_conf.cams.at(name);

    if (setup.axes.x.weight != 0) {
      double const beta = conf.axes.x.beta;
      double const disp = conf.axes.x.dispersion;
      double const size = std::sqrt(beta * em + std::pow(disp * es, 2));

      realizations.push_back(AxisRealization{
          .config = conf.axes.x,
          .setup = setup.axes.x,
          .measurement = size,
      });
    }

    if (setup.axes.z.weight != 0) {
      double const beta = conf.axes.z.beta;
      double const disp = conf.axes.z.dispersion;
      double const size = std::sqrt(beta * em + std::pow(disp * es, 2));

      realizations.push_back(AxisRealization{
          .config = conf.axes.z,
          .setup = setup.axes.z,
          .measurement = size,
      });
    }
  }

  double x[] = {1.0, 1.0};

  ASSERT_TRUE(v2k::ems::Estimate(realizations, x));
  ASSERT_NEAR(x[0], em, 1e-7);
  ASSERT_NEAR(x[1], es, 1e-7);
}

} // namespace ems
} // namespace v2k

int main(int argc, char **argv) {
  spdlog::set_level(spdlog::level::debug);
  testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
