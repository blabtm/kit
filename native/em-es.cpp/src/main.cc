#include <iostream>
#include <v2k/ccd.h>

int main(int argc, char **argv) {
  (void)argc;
  (void)argv;

  v2k::ccd::Config config("/workspaces/v2k/config/ccd/config.yaml");

  for (auto const &[name, conf] : config.cams) {
    std::cout << name << std::endl;
    std::cout << " x" << std::endl;
    std::cout << "  beta " << conf.axes.x.beta << std::endl;
    std::cout << "  disp " << conf.axes.x.dispersion << std::endl;
    std::cout << " z" << std::endl;
    std::cout << "  beta " << conf.axes.z.beta << std::endl;
    std::cout << "  disp " << conf.axes.z.dispersion << std::endl;
  }

  return 0;
}
