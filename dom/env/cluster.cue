package env

import (
	"net"
  "list"
)

// Node is an addressable Linux machine.
#Node: {
	name!: =~"^[a-z]+[a-z0-9]*$"

	// host is a node's LAN address.
	host!: string & net.IPv4

	// port is a node's SSH port.
	port: >0 & <100000 | *22
}

// Cluster is a collection of addressable unique Linux machines.
#Cluster: [#Node, ...#Node]

#ClusterValidator: {
  #cluster: #Cluster

  _clusterNameCheck: [for node in #cluster { "\(node.name)" }] & list.UniqueItems
  _clusterHostCheck: [for node in #cluster { "\(node.host)" }] & list.UniqueItems

  #AllowedNode: or([for node in #cluster {node.name}]) 
}
