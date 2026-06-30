package sm

import (
	"github.com/blabtm/v2k/dom/platform"
)

#Config: platform.#Config & {
  active: true
  deployment: {
    node: "master"
  }
}
