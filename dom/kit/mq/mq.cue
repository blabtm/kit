package mq

import (
	"list"

	"lab.kit/v1/dom/kit/svc"
)

#Listener: {
	type:    "tcp" | "sysinfo" | "ws"
	id:      string
	address: =~":[0-9]+"
}

#Config: svc.#Config & {
	enable: true
	listeners: [...#Listener] | *[
		{
			type:    "tcp"
			id:      "tcp"
			address: ":1883"
		},
		{
			type:    "ws"
			id:      "ws"
			address: ":1882"
		},
	]

	_listenersIdCheck: [for l in listeners {"\(l.id)"}] & list.UniqueItems
}
