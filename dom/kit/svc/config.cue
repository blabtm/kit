package svc

// Deployment specifies service's location
// in a locally-distributed cluster.
#Deployment: {
  #AllowedNode: string

	node?: [...#AllowedNode]

  ...
}

// Config is a desired configuration of the service.
#Config: {
	enable:     bool | *false
  version:    #Version
	deployment: #Deployment

	...
}
