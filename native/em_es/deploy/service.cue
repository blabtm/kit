package em_es

import v2k "github.com/blabtm/v2k/model"

v2k.#Service & {
	driver: "swarm"
	type:   "service"
	capabilities: {
		"live": true
	}
}
