package rp

import (
	"lab.kit/v1/dom/kit/env"
	"lab.kit/v1/dom/kit/svc"
)

#Secret: {
	user: string
	pass: string
}

#Config: svc.#Config & {
	#cluster: env.#Cluster

	secret: #Secret

	enable: true
	deployment: svc.#Deployment & {
		node:  [string, ...string]
		port:  int | *9092
	}

	superusers: [secret.user]
	auto_create_topics_enabled: bool | *true

	console: #ConsoleConfig & {
		#kafka: {
			for node in #cluster {
				if node.name == deployment.node[0] {
					host: node.host
				}
			}
			port: deployment.port
		}
		#secret: secret
	}
}
